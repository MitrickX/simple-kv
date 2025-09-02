package db

import (
	"fmt"

	"github.com/MitrickX/simple-kv/internal/interpreter"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/storage"
)

type DB struct {
	interpreter interpreter.Interpreter
	storage     storage.Storage
}

func NewDB(
	interpreter interpreter.Interpreter,
	storage storage.Storage,
) *DB {
	return &DB{
		interpreter: interpreter,
		storage:     storage,
	}
}

func (db *DB) Exec(query string) (string, error) {
	result, err := db.interpreter.Interpret(query)
	if err != nil {
		return "", fmt.Errorf("db exec fail: %w", err)
	}

	switch result.Command.CommandType {
	case parser.SetCommandType:
		db.storage.Set(result.Command.Arguments[0], result.Command.Arguments[1])
		return "ok", nil
	case parser.GetCommandType:
		val, exists := db.storage.Get(result.Command.Arguments[0])
		if exists {
			return fmt.Sprintf("val: %s", val), nil
		} else {
			return "none", nil
		}
	case parser.DelCommandType:
		db.storage.Del(result.Command.Arguments[0])
		return "ok", nil
	default:
		return "none", nil
	}
}
