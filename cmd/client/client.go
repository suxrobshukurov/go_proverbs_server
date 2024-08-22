package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	addr       = "localhost:12345"
	proto      = "tcp4"
	retryCount = 5               // number of retries
	retryDelay = 2 * time.Second // delay between retries
	readTimout = 5 * time.Second // timeout for read operations
)

func main() {

	var conn net.Conn
	var err error

	for i := 0; i < retryCount; i++ {
		conn, err = net.Dial(proto, addr)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to server (attempt %d/%d): %v", i+1, retryCount, err)
		time.Sleep(retryDelay)
	}
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to server", addr)

	reader := bufio.NewReader(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(readTimout))
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalf("Error reading from server: %v", err)
		}
		fmt.Print(string(b))
	}

}
