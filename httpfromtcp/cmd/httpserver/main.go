package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/request"
	"github.com/shubh-man007/BeGo/httpfromtcp/internal/response"
	"github.com/shubh-man007/BeGo/httpfromtcp/internal/server"
)

const port = 42069
const CRLF = "\r\n"

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

	body := []byte(res200)
	stat := response.StatusOK

	if path == "/yourproblem" {
		body = []byte(res400)
		stat = response.StatusBadRequest
	} else if path == "/myproblem" {
		body = []byte(res500)
		stat = response.StatusInternalServerError
	} else if path == "/video" {
		file, err := os.Open("assets/warehouse.mp4") // Set file name accordingly.
		if err != nil {
			body = []byte(res500)
			stat = response.StatusInternalServerError
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			body = []byte(res500)
			stat = response.StatusInternalServerError
		}

		h := response.GetDefaultHeaders(int(info.Size()))
		h.Replace("Content-Type", "video/mp4")

		w.WriteStatusLine(response.StatusOK)
		w.WriteHeaders(h)

		buffer := make([]byte, 32*1024) // 32KB chunks
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				w.WriteBody(buffer[:n])
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				body = []byte(res500)
				stat = response.StatusInternalServerError
			}
		}
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/stream") {
		target := req.RequestLine.RequestTarget
		res, err := http.Get("https://httpbin.org" + strings.TrimPrefix(target, "/httpbin"))
		if err != nil {
			body = []byte(res500)
			stat = response.StatusInternalServerError
		} else {
			w.WriteStatusLine(response.StatusOK)

			h.Delete("Content-Length")
			h.Set("transfer-encoding", "chunked")
			h.Replace("Content-Type", "text/plain")
			h.Set("Trailer", "X-Content-SHA256")
			h.Set("Trailer", "X-Content-Length")

			w.WriteHeaders(h)
			w.WriteTrailers(h)

			var fullBody []byte

			fmt.Printf(">> Proxy Response: \n")
			for {
				data := make([]byte, 32)
				n, err := res.Body.Read(data)
				if err != nil {
					break
				}
				w.WriteBody([]byte(fmt.Sprintf("%x%s", n, CRLF)))
				w.WriteBody(data[:n])
				w.WriteBody([]byte(CRLF))

				fullBody = append(fullBody, data[:n]...)

				log.Printf("\n%s\n", string(data[:n]))
			}
			w.WriteBody([]byte(fmt.Sprintf("0%s%s", CRLF, CRLF)))

			hash := sha256.Sum256(fullBody)
			h.Set("X-Content-SHA256", hex.EncodeToString(hash[:]))
			h.Set("X-Content-Length", strconv.Itoa(len(fullBody)))

			w.WriteTrailers(h)

			log.Println("\n-----------------")
			return
		}
	}

	// switch path {
	// case "/yourproblem":
	// 	body = []byte(res400)
	// 	stat = response.StatusBadRequest

	// case "/myproblem":
	// 	body = []byte(res500)
	// 	stat = response.StatusInternalServerError

	// default:
	// 	body = []byte(res200)
	// 	stat = response.StatusOK
	// }

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
