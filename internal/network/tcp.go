package network

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/interpreter"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/storage"
	"go.uber.org/zap"
)

const (
	MessageHello = "HELLO"
	MessageHi    = "HI"
	MessageBye   = "BYE"

	startBufSize = 4096
)

type TcpServer struct {
	config      *config.Config
	interpreter interpreter.Interpreter
	storage     storage.Storage
	logger      *zap.Logger
	connLimiter chan struct{}
}

func NewTcpServer(
	config *config.Config,
	interpreter interpreter.Interpreter,
	storage storage.Storage,
	logger *zap.Logger,
) *TcpServer {
	return &TcpServer{
		config:      config,
		interpreter: interpreter,
		storage:     storage,
		logger:      logger,
		connLimiter: make(chan struct{}, config.Network.MaxConnections),
	}
}

func (s *TcpServer) Start(ctx context.Context) error {
	addr := s.config.Network.Address
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		s.logger.Error("failed to listen", zap.String("address", addr), zap.Error(err))
		return err
	}
	defer ln.Close()

	s.logger.Info("tpc server listening",
		zap.String("address", ln.Addr().String()),
	)

	for {
		// limit number of connections using cond var
		s.connLimiter <- struct{}{}
		conn, err := ln.Accept()
		if err != nil {
			<-s.connLimiter
			s.logger.Error("failed to accept connection", zap.Error(err))
			continue
		}

		s.logger.Info("accepted connection", zap.String("remote", conn.RemoteAddr().String()), zap.Int("connCount", len(s.connLimiter)))
		go s.handleConn(conn)
	}
}

func (s *TcpServer) handleConn(conn net.Conn) {
	defer func() {
		<-s.connLimiter
		conn.Write([]byte(MessageBye))
		conn.Close()

		s.logger.Info("closed connection", zap.String("remote", conn.RemoteAddr().String()))

		if r := recover(); r != nil {
			s.logger.Error("panic happened", zap.Any("recovery", r))
		}
	}()

	s.handshake(conn)

	// set idle deadline
	conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.Network.IdleTimeout)))

	// init scanner with token size limit
	scanner := bufio.NewScanner(conn)
	bufSize := min(int(s.config.Network.MaxMessageSize), startBufSize)
	scanner.Buffer(make([]byte, bufSize), int(s.config.Network.MaxMessageSize))

	for {
		for scanner.Scan() {
			// move idle deadline
			conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.Network.IdleTimeout)))

			query := scanner.Text()

			s.logger.Debug("input query", zap.String("query", query))

			result, err := s.exec(query)

			s.logger.Debug("execute query", zap.String("result", result), zap.Error(err))

			if err != nil {
				fmt.Fprintf(conn, "%s\n", err.Error())
				continue
			}

			fmt.Fprintf(conn, "%s\n", result)
		}

		err := scanner.Err()
		if err != nil {
			if errors.Is(err, bufio.ErrTooLong) {
				fmt.Fprintf(conn, "error: %s\n", err.Error())
			}
			s.logger.Error("connection error", zap.Error(err))
		}

		break
	}

}

func (s *TcpServer) exec(query string) (string, error) {
	result, err := s.interpreter.Interpret(query)
	if err != nil {
		return "", fmt.Errorf("db exec fail: %w", err)
	}

	res := s.storage.Exec(&result.Command)
	if res.Err != nil {
		return "", res.Err
	}

	if result.Command.CommandType == parser.GetCommandType {
		if res.Ok {
			return fmt.Sprintf("val: %s", res.Val), nil
		} else {
			return "none", nil
		}
	}

	if res.Ok {
		return "ok", nil
	}

	return "none", nil
}

func (s *TcpServer) handshake(conn net.Conn) bool {
	var buf = make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		s.logger.Error("connection error cause of fail handshake, can't read hello message", zap.Error(err))
		return false
	}

	if string(buf[0:n]) != MessageHello {
		s.logger.Error("connection error cause of fail handshake, expect HELLO", zap.Binary("buf", buf))
		return false
	}

	n, err = conn.Write([]byte(MessageHi))
	if err != nil {
		s.logger.Error("connection error cause of fail handshake, can't send hi message", zap.Error(err))
		return false
	}

	if n != 2 {
		s.logger.Error("connection error cause of fail handshake, can't send fully hi message", zap.Int("n", n))
		return false
	}

	return true
}
