package interpreter

import (
	"fmt"

	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
)

type Interpreter interface {
	Interpret(query string) (*Result, error)
}

type interpreter struct {
	parser parser.Parser
}

func NewInterpreter(parser parser.Parser) Interpreter {
	return &interpreter{parser: parser}
}

func (i *interpreter) Interpret(query string) (*Result, error) {
	cmd, err := i.parser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("interpreter fails: %w", err)
	}

	return &Result{Command: *cmd}, nil
}
