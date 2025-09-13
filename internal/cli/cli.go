package cli

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	readWriteConnDeadlineTimeout = time.Second

	MessageHello = "HELLO"
	MessageHi    = "HI"
	MessageBye   = "BYE"
)

type Cli struct {
	input     io.Reader
	output    io.Writer
	errOutput io.Writer
	conn      net.Conn
}

func NewCli(
	input io.Reader,
	output io.Writer,
	errOutput io.Writer,
	conn net.Conn,
) *Cli {
	return &Cli{
		input:     input,
		output:    output,
		errOutput: errOutput,
		conn:      conn,
	}
}

func (c *Cli) handlshake() {
	c.conn.SetDeadline(time.Now().Add(readWriteConnDeadlineTimeout))

	buf := []byte(MessageHello)
	n, err := c.conn.Write(buf)
	if err != nil && n != 6 {
		_ = fmt.Sprintln(c.output, "\nSession ended cause of fail handshake.")
		c.conn.Close()
		os.Exit(0)
	}

	n, err = c.conn.Read(buf)
	if err != nil && string(buf[0:n]) != MessageHi {
		_ = fmt.Sprintln(c.output, "\nSession ended cause of fail handshake.")
		c.conn.Close()
		os.Exit(0)
	}

	c.conn.SetDeadline(time.Time{})
}

func (c *Cli) Go() {
	defer c.conn.Close()

	c.handlshake()

	go func() {
		scanner := bufio.NewScanner(c.input)
		for {
			fmt.Fprint(c.output, "> ")
			if !scanner.Scan() {
				break
			}
			text := scanner.Text()

			if _, err := fmt.Fprintf(c.conn, "%s\n", text); err != nil {
				fmt.Fprintf(c.errOutput, "failed to send: %v\n", err)
				break
			}
		}
	}()

	for {
		serverReader := bufio.NewScanner(c.conn)
		if serverReader.Scan() {
			text := serverReader.Text()
			if text == MessageBye {
				fmt.Fprintln(c.errOutput, "server send bye message and closed connection")
				return
			}
			fmt.Fprintln(c.output, text)
		} else {
			fmt.Fprintln(c.errOutput, "server closed connection")
		}
	}
}
