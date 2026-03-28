package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go run main.go [--timeout=10s] host port")
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to", address)

	var wg sync.WaitGroup
	wg.Add(2)

	// 1. socket → stdout
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, conn)
	}()

	// 2. stdin → socket
	go func() {
		defer wg.Done()

		reader := bufio.NewReader(os.Stdin)
		for {
			data, err := reader.ReadBytes('\n')

			if len(data) > 0 {
				_, writeErr := conn.Write(data)
				if writeErr != nil {
					return
				}
			}

			if err != nil {
				if err == io.EOF {
					// Ctrl+D
					conn.Close()
				}
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println("\nConnection closed")
}