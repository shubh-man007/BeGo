package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/request"
	"github.com/shubh-man007/BeGo/httpfromtcp/internal/response"
)

// const resp = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"

type Server struct {
	Port     int
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

func NewServer() *Server {
	return &Server{}
}

func NewHandlerError() *HandlerError {
	return &HandlerError{}
}

type Handler func(w *response.Writer, req *request.Request)

func (herr *HandlerError) WriteErrorResponse(w io.Writer) error {
	code := herr.Code
	message := herr.Message

	_, err := response.WriteStatusLine(w, code)
	if err != nil {
		return errors.New("could not write error code to connection")
	}

	h := response.GetDefaultHeaders(len(message))
	err = response.WriteResHeaders(w, h)
	if err != nil {
		return errors.New("could not write error headers to connection")
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return errors.New("could not write error message to connection")
	}

	return nil
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		herr := &HandlerError{
			Code:    response.StatusBadRequest,
			Message: err.Error(),
		}
		herr.WriteErrorResponse(conn)
		return
	}

	w := response.NewWriter(conn)
	s.handler(w, req)
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

func Serve(port int, h Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return &Server{}, fmt.Errorf("failed to listen at address %v : %s", port, err.Error())
	}

	s := NewServer()
	s.Port = port
	s.handler = h
	s.listener = listener
	s.closed.Store(false)

	go s.listen()

	return s, nil
}
