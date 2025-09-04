package main

import (
	"fmt"
	"os"

	"github.com/MitrickX/simple-kv/internal/cli"
	"github.com/MitrickX/simple-kv/internal/db"
	"github.com/MitrickX/simple-kv/internal/interpreter"
	"github.com/MitrickX/simple-kv/internal/interpreter/parser"
	"github.com/MitrickX/simple-kv/internal/storage"
	"github.com/MitrickX/simple-kv/internal/storage/engine"
)

func main() {
	fmt.Println("Support commands: SET/GET/DEL")
	fmt.Println("SET key value")
	fmt.Println("GET key")
	fmt.Println("DEL key")

	parser := parser.NewParser()
	interpreter := interpreter.NewInterpreter(parser)
	engine := engine.NewEngine()
	storage := storage.NewStorage(engine)
	db := db.NewDB(interpreter, storage)
	cli.NewCli(os.Stdin, os.Stdout, os.Stderr, db).Go()
}
