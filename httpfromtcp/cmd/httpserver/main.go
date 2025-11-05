package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/request"
	"github.com/shubh-man007/BeGo/httpfromtcp/internal/response"
	"github.com/shubh-man007/BeGo/httpfromtcp/internal/server"
)

const port = 42069

// Client template:
const res400 = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`

const res500 = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`

const res200 = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

func RequestPath(w *response.Writer, req *request.Request) {
	path := req.RequestLine.RequestTarget
	h := response.GetDefaultHeaders(0)

	var body []byte
	var stat response.StatusCode

	switch path {
	case "/yourproblem":
		body = []byte(res400)
		stat = response.StatusBadRequest

	case "/myproblem":
		body = []byte(res500)
		stat = response.StatusInternalServerError

	default:
		body = []byte(res200)
		stat = response.StatusOK
	}

	h.Replace("Content-Length", strconv.Itoa(len(body)))
	h.Replace("Content-Type", "text/html")

	if err := w.WriteStatusLine(stat); err != nil {
		log.Printf("Error writing status line: %v", err)
		return
	}

	if err := w.WriteHeaders(h); err != nil {
		log.Printf("Error writing headers: %v", err)
		return
	}

	if _, err := w.WriteBody(body); err != nil {
		log.Printf("Error writing body: %v", err)
		return
	}

	res := w.LogResponse(stat, h, string(body))
	log.Printf("\nResponse: \n%s\n", res)
}

func main() {
	server, err := server.Serve(port, RequestPath)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("\nServer gracefully stopped")
}
