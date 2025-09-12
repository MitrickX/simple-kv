package cli

import (
	"bufio"
	"fmt"
	"io"
)

type Cli struct {
	input     io.Reader
	output    io.Writer
	errOutput io.Writer
	conn      io.ReadWriter
}

func NewCli(
	input io.Reader,
	output io.Writer,
	errOutput io.Writer,
	conn io.ReadWriter,
) *Cli {
	return &Cli{
		input:     input,
		output:    output,
		errOutput: errOutput,
		conn:      conn,
	}
}

func (cli *Cli) Go() {
	scanner := bufio.NewScanner(cli.input)
	serverReader := bufio.NewScanner(cli.conn)
	for {
		fmt.Fprint(cli.output, "> ")
		if !scanner.Scan() {
			break
		}
		text := scanner.Text()
		if _, err := fmt.Fprintf(cli.conn, "%s\n", text); err != nil {
			fmt.Fprintf(cli.errOutput, "failed to send: %v\n", err)
			break
		}
		if serverReader.Scan() {
			fmt.Fprintln(cli.output, serverReader.Text())
		} else {
			fmt.Fprintln(cli.errOutput, "server closed connection")
			break
		}
	}
}
