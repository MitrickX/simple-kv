package network

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/db"
	"go.uber.org/zap"
)

type TcpServer struct {
	config      *config.Config
	db          *db.DB
	logger      *zap.Logger
	connLimiter chan struct{}
}

func NewTcpServer(
	config *config.Config,
	db *db.DB,
	logger *zap.Logger,
) *TcpServer {
	return &TcpServer{
		config:      config,
		db:          db,
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

	s.logger.Info("tpc server listening", zap.String("address", ln.Addr().String()))

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
		conn.Close()

		s.logger.Info("closed connection", zap.String("remote", conn.RemoteAddr().String()))

		if r := recover(); r != nil {
			s.logger.Error("panic happened", zap.Any("recovery", r))
		}
	}()

	s.handshake(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		query := scanner.Text()

		s.logger.Debug("input query", zap.String("query", query))

		result, err := s.db.Exec(query)

		s.logger.Debug("execute query", zap.String("result", result), zap.Error(err))

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err.Error())
			continue
		}

		fmt.Fprintf(conn, "%s\n", result)
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("connection error", zap.Error(err))
	}
}

func (s *TcpServer) handshake(conn net.Conn) bool {
	var buf = make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		s.logger.Error("connection error cause of fail handshake, can't read hello message", zap.Error(err))
		return false
	}

	if string(buf[0:n]) != "HELLO" {
		s.logger.Error("connection error cause of fail handshake, expect HELLO", zap.Binary("buf", buf))
		return false
	}

	n, err = conn.Write([]byte("HI"))
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
