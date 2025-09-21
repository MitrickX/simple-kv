package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/storage/engine"
	"github.com/MitrickX/simple-kv/internal/wal"
)

type Result struct {
	Ok  bool
	Val string
	Err error
}

type Storage interface {
	Exec(cmd *parser.Command) Result
	Run(ctx context.Context)
}

func NewStorage(cfg *config.Config, engine engine.Engine, wal wal.WAL) Storage {
	return &storage{
		cfg:    cfg,
		engine: engine,
		wal:    wal,
		mx:     &sync.RWMutex{},
	}
}

type storage struct {
	cfg    *config.Config
	engine engine.Engine
	wal    wal.WAL
	mx     *sync.RWMutex
}

func (s *storage) Run(ctx context.Context) {
	go func() {
		fmt.Println("ticket", time.Duration(s.cfg.WAL.FlushingBatchTimeout))
		ticker := time.NewTicker(time.Duration(s.cfg.WAL.FlushingBatchTimeout))
		defer ticker.Stop()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.walFlush()
		}
	}()
}

func (s *storage) Exec(cmd *parser.Command) Result {
	if cmd.CommandType != parser.SetCommandType &&
		cmd.CommandType != parser.DelCommandType &&
		cmd.CommandType != parser.GetCommandType {
		return Result{
			Err: fmt.Errorf("unknown command type: %v", cmd.CommandType),
		}
	}

	if cmd.CommandType == parser.GetCommandType {
		s.mx.RLock()
		defer s.mx.RUnlock()
		val, ok := s.engine.Get(cmd.Arguments[0])
		return Result{
			Ok:  ok,
			Val: val,
		}
	}

	s.mx.Lock()
	defer s.mx.Unlock()
	err := s.wal.Write(string(cmd.CommandType))
	fmt.Println("WAL-write", err)
	if err != nil {
		return Result{
			Err: err,
		}
	}

	if cmd.CommandType == parser.SetCommandType {
		s.engine.Set(cmd.Arguments[0], cmd.Arguments[1])
	}

	s.engine.Del(cmd.Arguments[0])

	return Result{
		Ok: true,
	}
}

func (s *storage) walFlush() error {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.wal.Flush()
}
