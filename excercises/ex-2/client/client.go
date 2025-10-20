package client

import (
	"bufio"
	"log"
	"net"
	"os"
)

// Server:  Listen, Listen.Accept, Conn Interface (Read, Write, Close, Addr, etc.)
// Client: Dial, ResolveAddr, Conn

func ConnClient() {
	const loggerName = "[CLIENT]"
	const port = ":8080"

	log.SetPrefix(loggerName + "\t")

	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatalf("Failed to resolve address at port %s : %v", port, err)
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to connect to server at %s", conn.RemoteAddr())
	}
	defer conn.Close()

	log.Printf("Connected to server at %s", conn.RemoteAddr())

	rd := bufio.NewReader(os.Stdin)
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			return
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Printf("Error writing to server: %v", err)
			return
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from server: %v", err)
			return
		}
		log.Printf("Server: %s", string(buf[:n]))
	}
}
