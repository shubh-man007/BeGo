package server

import (
	"bufio"
	"io"
	"log"
	"net"
)

// Server:  Listen, Listen.Accept, Conn Interface (Read, Write, Close, Addr, etc.)
// Client: Dial, ResolveAddr, Conn

func RunServer() {
	// Server:
	const loggerName = "[SERVER]"
	const port = ":8080"

	log.SetPrefix(loggerName + "\t")

	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatalf("Failed to resolve address at port %s : %v", port, err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen at address %v : %s", addr, err.Error())
	}
	defer listener.Close()

	log.Printf("Listening at port %s", listener.Addr())

	for {
		// Ideally each connection should be handled via its own goroutine
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept request : %s", err.Error())
			continue
		}
		log.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		go func(c net.Conn) {
			defer c.Close()

			rd := bufio.NewReader(conn)
			for {
				str, err := rd.ReadString('\n')
				if err == io.EOF {
					// client closed connection
					log.Printf("[%s]\t Disconnected", c.RemoteAddr())
					return
				}
				log.Printf("[%s]\t %s", c.RemoteAddr(), str)
				c.Write([]byte("Echo: " + str))
			}
		}(conn)
	}
}
