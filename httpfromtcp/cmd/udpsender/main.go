package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const port = ":42069"
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("Failed to resolve address %s: %s", port, err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to connect to port %s: %s", port, err.Error())
	}

	defer conn.Close()

	rd := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(">")
		input, _ := rd.ReadString('\n')
		conn.Write([]byte(input))
	}
}
