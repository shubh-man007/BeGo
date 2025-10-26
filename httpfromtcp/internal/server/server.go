package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

const resp = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"

type Server struct {
	Port     int
	listener net.Listener
	closed   atomic.Bool
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		s.closed.Store(false)
		return errors.New("could not close listener")
	}
	s.closed.Store(true)
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	_, err := conn.Write([]byte(resp))
	if err != nil {
		log.Printf("Failed to write back to client")
	}
	log.Printf("Server: \n%s", resp)
}

func (s *Server) listen() {
	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Failed to accept request: %s", err.Error())
			continue
		}
		log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return &Server{}, fmt.Errorf("failed to listen at address %v : %s", port, err.Error())
	}

	s := NewServer()
	s.Port = port
	s.listener = listener
	s.closed.Store(false)

	go s.listen()

	return s, nil
}
