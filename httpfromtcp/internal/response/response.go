package response

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/headers"
)

const CRLF = "\r\n"

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
	StatusUnrecog             StatusCode = -1
)

const (
	WriterStatusInit = iota
	WriterStatusHeader
	WriterStatusBody
	WriterStatusDone
)

type Writer struct {
	Status int
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) (StatusCode, error) {
	var rp string
	switch statusCode {
	case StatusOK:
		rp = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		rp = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		rp = "HTTP/1.1 500 Internal Server Error"
	default:
		return StatusUnrecog, errors.New("unsupported status code")
	}

	_, err := w.Write([]byte(rp + CRLF))
	if err != nil {
		return statusCode, errors.New("could not write to connection")
	}
	return statusCode, nil
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteResHeaders(w io.Writer, h *headers.Headers) error {
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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.Status != WriterStatusInit {
		return errors.New("state mismatch, status line parsed or skipped")
	}

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

	_, err := w.writer.Write([]byte(rp + CRLF))
	if err != nil {
		return err
	}

	w.Status = WriterStatusHeader

	return nil
}

func (w *Writer) WriteHeaders(headers *headers.Headers) error {
	if w.Status != WriterStatusHeader {
		return errors.New("state mismatch, headers parsed or skipped")
	}

	for key, value := range headers.Iter() {
		fieldLine := fmt.Sprintf("%s: %s%s", key, value, CRLF)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(CRLF))
	if err != nil {
		return err
	}

	w.Status = WriterStatusBody

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.Status != WriterStatusBody && w.Status != WriterStatusDone {
		return 0, errors.New("state mismatch: must write headers before body")
	}

	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}

	w.Status = WriterStatusDone
	return n, nil
}

func (w *Writer) WriteTrailers(headers *headers.Headers) error {
	if w.Status != WriterStatusDone {
		return errors.New("state mismatch: must write headers before body")
	}

	trailerStr := headers.Get("Trailer")
	trailers := strings.Split(trailerStr, ",")

	for _, key := range trailers {
		value := headers.Get(key)
		fieldLine := fmt.Sprintf("%s: %s%s", key, value, CRLF)
		_, err := w.writer.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}

	_, err := w.writer.Write([]byte(CRLF))
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) LogResponse(statusCode StatusCode, h *headers.Headers, body string) string {
	if w.Status != WriterStatusDone {
		return errors.New("state mismatch, response not formed yet").Error()
	}

	var res string
	var rp string
	switch statusCode {
	case StatusOK:
		rp = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		rp = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		rp = "HTTP/1.1 500 Internal Server Error"
	default:
		return errors.New("unsupported status code").Error()
	}

	rp += CRLF
	res += rp

	for key, value := range h.Iter() {
		fieldLine := fmt.Sprintf("%s: %s%s", key, value, CRLF)
		res += fieldLine
	}

	res += CRLF
	res += body

	return res
}
