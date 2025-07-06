package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func echoUpper(w io.Writer, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		newText := strings.ToUpper(scanner.Text())
		fmt.Fprintf(w, "%s \n", newText)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error: %s", err)
	}
}

func main() {
	const name = "uppertcp"
	log.SetPrefix(name + "\t")

	port := flag.Int("port", 8080, "Port to listen to")

	flag.Parse()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatalf("Failed to listen to port %d: %v", listener.Addr(), err)
	}

	defer listener.Close()

	log.Printf("listening at localhost: %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection request at port %d: %v", listener.Addr(), err)
		}

		go echoUpper(conn, conn)
	}
}
