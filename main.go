package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

var (
	counter int

	listenAddr = "localhost:8080"

	server = []string{
		"localhost:5001",
		"localhost:5002",
		"localhost:5003",
	}
)

func main() {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listen failure: %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err)
		}

		server := chooseServer()
		fmt.Println("counter=%d and server=%s\n", counter, server)
		go func() {
			proxy(server, conn)
		}()
	}
}

func proxy(server string, c net.Conn) {
	sc, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("failed to connect to server %s: %v", server, err)
	}

	go io.Copy(sc, c)
	go io.Copy(c, sc)
}

func chooseServer() string {
	s := server[counter%len(server)]
	counter++

	return s
}
