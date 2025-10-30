package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/headers"
)

const CRLF = "\r\n"

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var rp string
	switch statusCode {
	case StatusOK:
		rp = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		rp = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		rp = "HTTP/1.1 500 Internal Server Error"
	default:
		return errors.New("unsupported status code")
	}

	_, err := w.Write([]byte(rp + CRLF))
	if err != nil {
		return errors.New("could not write to connection")
	}
	return nil
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, h *headers.Headers) error {
	for key, value := range h.Iter() {
		fieldLine := fmt.Sprintf("%s: %s%s", key, value, CRLF)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return errors.New("could not write headers to connection")
		}
	}
	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return errors.New("could not write final CRLF")
	}
	return nil
}
