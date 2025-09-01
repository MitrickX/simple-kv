package parser

type CommandType string

const (
    SetCommandType CommandType = "SET"
    GetCommandType CommandType = "GET"
    DelCommandType CommandType = "DEL"
)

type Command struct {
    CommandType CommandType
    Arguments []string
}