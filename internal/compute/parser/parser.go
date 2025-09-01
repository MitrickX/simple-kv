package parser

import (
	"strings"
)

type Parser interface {
	Parse(string) (*Command, error)
}

func NewParser() Parser {
	return &parser{}
}

type parser struct {}

func (p *parser) Parse(query string) (*Command, error) {
	parts := strings.Split(query, " ")

	tokens := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			tokens = append(tokens, part)
		}
	}

	if len(tokens) == 0 {
		return nil, ErrNoTokensInQuery
	}

	switch CommandType(tokens[0]) {
	case SetCommandType:
		if len(tokens) < 3 {
			return nil, ErrNoEnoughArgumentsForSetCommand
		}
		return &Command{
			CommandType: SetCommandType,
			Arguments:   tokens[1:3],
		}, nil
	case GetCommandType:
		if len(tokens) < 2 {
			return nil, ErrNoEnoughArgumentsForGetCommand
		}
		return &Command{
			CommandType: GetCommandType,
			Arguments:   tokens[1:2],
		}, nil
	case DelCommandType:
		if len(tokens) < 2 {
			return nil, ErrNoEnoughArgumentsForDelCommand
		}
		return &Command{
			CommandType: DelCommandType,
			Arguments:   tokens[1:2],
		}, nil
	default:
        return nil, ErrUnknownCommandType
	}
}
