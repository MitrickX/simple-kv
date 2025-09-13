package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MitrickX/simple-kv/internal/cli"
)

const (
	tcpDialTimeout = time.Second
)

func main() {
	address := flag.String("address", "", "network address to connect to TCP server")
	flag.Parse()

	if *address == "" {
		fmt.Println("Usage: cli --address <address>")
		os.Exit(1)
	}

	fmt.Println("Support commands: SET/GET/DEL")
	fmt.Println("SET key value")
	fmt.Println("GET key")
	fmt.Println("DEL key")

	conn, err := net.DialTimeout("tcp", *address, tcpDialTimeout)
	if err != nil {
		fmt.Printf("failed to connect to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", *address)

	// Handle Ctrl+C (SIGINT) to exit gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	go func() {
		<-sigCh
		fmt.Println("\nSession ended.")
		conn.Close()
		os.Exit(0)
	}()

	cli := cli.NewCli(os.Stdin, os.Stdout, os.Stderr, conn)
	cli.Go()
}
