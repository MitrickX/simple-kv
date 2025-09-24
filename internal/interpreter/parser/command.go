package parser

import "strings"

type CommandType string

const (
	SetCommandType CommandType = "SET"
	GetCommandType CommandType = "GET"
	DelCommandType CommandType = "DEL"
)

type Command struct {
	CommandType CommandType
	Arguments   []string
}

func (c Command) String() string {
	return string(c.CommandType) + " " + strings.Join(c.Arguments, " ")
}
