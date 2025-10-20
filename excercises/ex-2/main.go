package main

import (
	"flag"
	"log"

	"github.com/shubh-man007/BeGo/httpfromtcp/excercises/ex-2/client"
	"github.com/shubh-man007/BeGo/httpfromtcp/excercises/ex-2/server"
)

func main() {
	mode := flag.String("mode", "server", "Client-Server switch")
	flag.Parse()

	switch *mode {
	case "server":
		log.Println("Starting server...")
		server.RunServer()
	case "client":
		log.Println("Starting client...")
		client.ConnClient()
	default:
		log.Fatal("Unknown mode. Use -mode=server or -mode=client")
	}
}
