package main

import (
	"context"
	"flag"
	"log"
	"os"

	"strings"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/interpreter"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/network"
	"github.com/MitrickX/simple-kv/internal/storage"
	"github.com/MitrickX/simple-kv/internal/storage/engine"
	utilsOs "github.com/MitrickX/simple-kv/internal/utils/os"
	utilsTime "github.com/MitrickX/simple-kv/internal/utils/time"
	"github.com/MitrickX/simple-kv/internal/wal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	cfg := config.Default()
	if *configPath != "" {
		var err error
		cfg, err = config.Parse(*configPath)
		if err != nil {
			log.Fatalf("failed to parse config: %v\n", err)
		}
	}

	logger := buildZap(&cfg)
	defer logger.Sync()

	logger.Info("config parsed", zap.Any("config", cfg))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parser := parser.NewParser()
	interpreter := interpreter.NewInterpreter(parser)
	engine := engine.NewEngine()
	wal := wal.NewWAL(&cfg, utilsOs.NewOS(), utilsTime.NewTime())

	storage := storage.NewStorage(&cfg, engine, wal)
	storage.Run(ctx)

	server := network.NewTcpServer(&cfg, interpreter, storage, logger)
	if err := server.Start(ctx); err != nil {
		logger.Fatal("server exited with error", zap.Error(err))
	}
}

func buildZap(cfg *config.Config) *zap.Logger {
	var level zapcore.Level
	switch strings.ToLower(cfg.Logging.Level) {
	case config.LoggingLevelDebug:
		level = zapcore.DebugLevel
	case config.LoggingLevelInfo:
		level = zapcore.InfoLevel
	case config.LoggingLevelWarning:
		level = zapcore.WarnLevel
	case config.LoggingLevelError:
		level = zapcore.ErrorLevel
	case config.LoggingLevelPanic:
		level = zapcore.PanicLevel
	case config.LoggingLevelFatal:
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.OutputPaths = []string{cfg.Logging.Output}

	logger, err := zapCfg.Build()
	if err != nil {
		log.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}
