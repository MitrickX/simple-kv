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

	res := db.storage.Exec(&result.Command)
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
