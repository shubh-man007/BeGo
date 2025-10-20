// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net"
// 	"strings"

// 	"github.com/spkg/bom"
// )

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lineChan := make(chan string)

// 	go func() {
// 		defer close(lineChan)
// 		defer f.Close()

// 		var tempLine string
// 		buf := make([]byte, 8)

// 		for {
// 			n, err := f.Read(buf)
// 			if n > 0 {
// 				tempLine += string(bom.Clean(buf[:n]))
// 				normTempLine := strings.ReplaceAll(tempLine, "\r\n", "\n")
// 				lines := strings.Split(normTempLine, "\n")

// 				for i := 0; i < len(lines)-1; i++ {
// 					lineChan <- lines[i]
// 				}

// 				tempLine = lines[len(lines)-1]
// 			}

// 			if err == io.EOF {
// 				if len(tempLine) > 0 {
// 					lineChan <- tempLine
// 				}
// 				break
// 			}

// 			if err != nil {
// 				log.Printf("error reading file: %v", err)
// 				break
// 			}
// 		}
// 	}()

// 	return lineChan
// }

// func readFromConn(c net.Conn) {
// 	lines := getLinesChannel(c)

// 	for v := range lines {
// 		fmt.Printf("%s\n", v)
// 	}

// 	log.Printf("Conn at %s, has been closed", c.RemoteAddr())
// }

// func main() {
// 	const port = ":42069"
// 	const name = "uppertcp"
// 	log.SetPrefix(name + "\t")

// 	listener, err := net.Listen("tcp", port)
// 	if err != nil {
// 		log.Fatalf("Failed to listen at port %s: %s\n", listener.Addr(), err.Error())
// 	}

// 	defer listener.Close()

// 	log.Printf("Listening at port %s", listener.Addr())

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Fatalf("Error in connecting at port %s: %s", listener.Addr(), err.Error())
// 		}

// 		log.Printf("Accepted connection request at port %s", conn.RemoteAddr())

// 		go readFromConn(conn)
// 	}
// }

// NOTE:
// Curl sends the header first and then the body, because of a missing CRLF it believes that there are more headers coming in and keeps the connection open, buffering the request payload.
// And when we interrupt the curl request, the data present in the buffer is sent to the server and is printed out in the console.

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/request"
)

func main() {
	const port = ":42069"
	const name = "uppertcp"
	log.SetPrefix(name + "\t")

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen at port %s: %s\n", port, err.Error())
	}
	defer listener.Close()

	log.Printf("Listening at %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("Error in parsing request: %s", err.Error())
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
