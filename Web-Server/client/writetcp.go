package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const name = "tcpwriter"
	log.SetPrefix(name + "\t")

	port := flag.Int("Port", 8080, "Port to attend to.")

	flag.Parse()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{Port: *port})
	if err != nil {
		log.Fatalf("Error in connection to port %d : %v", conn.RemoteAddr(), err)
	}

	defer conn.Close()

	go func() {
		// To read from the server
		for connScanner := bufio.NewScanner(conn); connScanner.Scan(); {
			fmt.Printf("SERVER: %s\n", connScanner.Text())
			if err := connScanner.Err(); err != nil {
				log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
			}
		}
	}()

	// To write to the server
	for stdinScanner := bufio.NewScanner(os.Stdin); stdinScanner.Scan(); {
		fmt.Printf("CLIENT: %s \n", stdinScanner.Text())
		if _, err := conn.Write(stdinScanner.Bytes()); err != nil {
			log.Fatalf("Could not write to server atr port %d: %v", conn.RemoteAddr(), err)
		}
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Fatalf("Could not write to server atr port %d: %v", conn.RemoteAddr(), err)
		}
		if err := stdinScanner.Err(); err != nil {
			log.Fatalf("error reading from %s: %v", conn.RemoteAddr(), err)
		}
	}
}
