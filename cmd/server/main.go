package main

import (
	"flag"
	"log"
	"os"

	"strings"

	"github.com/MitrickX/simple-kv/internal/config"
	"github.com/MitrickX/simple-kv/internal/db"
	"github.com/MitrickX/simple-kv/internal/interpreter"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/network"
	"github.com/MitrickX/simple-kv/internal/storage"
	"github.com/MitrickX/simple-kv/internal/storage/engine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		log.Fatalln("Usage: server --config <path>")
	}

	cfg, err := config.Parse(*configPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v\n", err)
	}

	logger := buildZap(&cfg)
	defer logger.Sync()

	parser := parser.NewParser()
	interpreter := interpreter.NewInterpreter(parser)
	engine := engine.NewEngine()
	storage := storage.NewStorage(engine)
	db := db.NewDB(interpreter, storage)

	server := network.NewTcpServer(&cfg, db, logger)
	if err := server.Start(); err != nil {
		logger.Fatal("server exited with error", zap.Error(err))
	}
}

func buildZap(cfg *config.Config) *zap.Logger {
	var level zapcore.Level
	switch strings.ToLower(cfg.Logging.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.OutputPaths = []string{cfg.Logging.Output, os.Stderr.Name()}

	logger, err := zapCfg.Build()
	if err != nil {
		log.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}
