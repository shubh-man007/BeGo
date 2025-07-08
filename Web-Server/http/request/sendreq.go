package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	method, host, path string
	port               int
)

func main() {
	flag.StringVar(&method, "method", "GET", "Type of HTTP request")
	flag.StringVar(&host, "host", "localhost", "Host to connect to")
	flag.IntVar(&port, "port", 8080, "Port to attend to")
	flag.StringVar(&path, "path", "/", "Path to fetch")

	flag.Parse()

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Unable to resolve address: %s:%d", host, port)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("Error in connecting to port %s : %v", conn.RemoteAddr(), err)
	}

	defer conn.Close()

	log.Printf("Connected to port %s from port %s", conn.RemoteAddr(), conn.LocalAddr())

	reqfields := []string{
		fmt.Sprintf("%s %s HTTP/1.1", method, path),
		"Host: " + host,
		"User-Agent: httpget",
		"",
	}

	req := strings.Join(reqfields, "\r\n") + "\r\n"

	conn.Write([]byte(req))
	log.Printf("sent request:\n%s", req)

	log.Printf("Reading from server at port %s", conn.RemoteAddr())
	for scanner := bufio.NewScanner(conn); scanner.Scan(); {
		line := scanner.Bytes()
		if _, err := fmt.Fprintf(os.Stdout, "%s\n", line); err != nil {
			log.Printf("error writing to connection: %s", err)
		}
		if scanner.Err() != nil {
			log.Printf("error reading from connection: %s", err)
			return
		}
	}

}
