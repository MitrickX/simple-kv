package parser

import "errors"

var (
	ErrNoTokensInQuery                = errors.New("parser error: no tokens in query")
	ErrNoEnoughArgumentsForSetCommand = errors.New("parser error: no enough arguments for set command")
	ErrNoEnoughArgumentsForGetCommand = errors.New("parser error: no enough arguments for get command")
	ErrNoEnoughArgumentsForDelCommand = errors.New("parser error: no enough arguments for del command")
	ErrUnknownCommandType             = errors.New("parser error: unknown command type")
)
