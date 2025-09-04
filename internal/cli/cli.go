package cli

import (
	"bufio"
	"io"

	"github.com/MitrickX/simple-kv/internal/db"
)

type Cli struct {
	scanner *bufio.Scanner
	db      *db.DB
	out     io.Writer
	errOut  io.Writer
}

func NewCli(
	reader io.Reader,
	out io.Writer,
	errOut io.Writer,
	db *db.DB,
) *Cli {
	scanner := bufio.NewScanner(reader)
	return &Cli{
		scanner: scanner,
		db:      db,
		out:     out,
		errOut:  errOut,
	}
}

func (cli *Cli) Go() {
	for cli.scanner.Scan() {
		query := cli.scanner.Text()
		result, err := cli.db.Exec(query)
		if err != nil {
			cli.errOut.Write([]byte(err.Error()))
			cli.out.Write([]byte{'\n'})
			continue
		}

		cli.out.Write([]byte(result))
		cli.out.Write([]byte{'\n'})
	}
}
