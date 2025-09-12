package network

import (
	"bufio"
	"fmt"
	"net"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/db"
	"go.uber.org/zap"
)

type TcpServer struct {
	config *config.Config
	db     *db.DB
	logger *zap.Logger
}

func NewTcpServer(
	config *config.Config,
	db *db.DB,
	logger *zap.Logger,
) *TcpServer {
	return &TcpServer{
		config: config,
		db:     db,
		logger: logger,
	}
}

func (s *TcpServer) Start() error {
	addr := s.config.Network.Address
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Error("failed to listen", zap.String("address", addr), zap.Error(err))
		return err
	}
	defer ln.Close()

	s.logger.Info("tpc server listening", zap.String("address", ln.Addr().String()))

	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error("failed to accept connection", zap.Error(err))
			continue
		}
		s.logger.Info("accepted connection", zap.String("remote", conn.RemoteAddr().String()))
		go s.handleConn(conn)
	}
}

func (s *TcpServer) handleConn(conn net.Conn) {
	defer conn.Close()
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
	} else {
		s.logger.Info("closed connection", zap.String("remote", conn.RemoteAddr().String()))
	}
}
