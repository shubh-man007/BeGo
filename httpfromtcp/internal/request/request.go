package request

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/shubh-man007/BeGo/httpfromtcp/internal/headers"
)

const CRLF = "\r\n"

const (
	stateInitialized = iota
	stateDone
	requestStateParsingHeaders
	requestStateParsingBody
)

type Request struct {
	RequestLine RequestLine
	state       int
	Header      *headers.Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func parseRequestLine(request string) (RequestLine, int, error) {
	idx := strings.Index(string(request), CRLF)
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	line := request[:idx]
	elements := strings.Fields(line)

	if len(elements) == 3 {
		// Validate HTTP method:
		if !IsUpper(elements[0]) {
			return RequestLine{}, idx, errors.New("invalid HTTP method, must be uppercase")
		}

		// Validate HTTP version
		if !strings.HasPrefix(elements[2], "HTTP/") {
			return RequestLine{}, idx, errors.New("invalid HTTP version format")
		}

		version := strings.TrimPrefix(elements[2], "HTTP/")
		if version != "1.1" {
			return RequestLine{}, idx, errors.New("unsupported HTTP version, only 1.1 is allowed")
		}

		reqStruct := RequestLine{}

		reqStruct.HttpVersion = version
		reqStruct.RequestTarget = elements[1]
		reqStruct.Method = elements[0]

		return reqStruct, idx, nil
	}

	return RequestLine{}, idx, errors.New("failed to parse request line: Incomplete request line")
}

func (r *Request) Parse(data []byte) (int, error) {
	switch r.state {
	case stateInitialized:
		reqLine, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		consumed := n + len("\r\n")
		r.RequestLine = reqLine
		r.state = requestStateParsingHeaders
		return consumed, nil

	case requestStateParsingHeaders:
		n, done, err := r.Header.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.state = requestStateParsingBody
		}

		return n, nil

	case requestStateParsingBody:
		CLVal := r.Header.Get("Content-Length")
		if CLVal == "" {
			r.state = stateDone
			return 0, nil
		}

		r.Body = append(r.Body, data...)

		CLInt, err := strconv.Atoi(CLVal)
		if err != nil {
			return 0, errors.New("error: invalid content length value")
		}

		if len(r.Body) > CLInt {
			return 0, errors.New("error: body length greater than header specified value")
		}

		// if len(r.Body) < CLInt {
		// 	return 0, errors.New("error: body length less than header specified value")
		// }

		if len(r.Body) == CLInt {
			r.state = stateDone
		}
		return len(data), nil

	case stateDone:
		return 0, errors.New("error: trying to read data in a done state")

	default:
		return 0, errors.New("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	const buffSize = 8
	buff := make([]byte, buffSize)
	readToIndex := 0 //bytes till which buff are filled

	req := &Request{
		state:  stateInitialized,
		Header: headers.NewHeaders(),
	}

	for req.state != stateDone {
		// parse whatever we already have.
		if readToIndex > 0 || req.state != stateInitialized {
			consumed, err := req.Parse(buff[:readToIndex])
			if err != nil {
				return nil, err
			}
			if consumed > 0 {
				copy(buff, buff[consumed:readToIndex])
				readToIndex -= consumed
				// attempt parsing again before reading more.
				continue
			}
			if req.state == stateDone {
				break
			}
		}

		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				// Final attempt to parse remaining buffered data on EOF
				if readToIndex > 0 || req.state != stateInitialized {
					if _, perr := req.Parse(buff[:readToIndex]); perr != nil {
						return nil, perr
					}
				}
				break
			}
			return nil, err
		}

		readToIndex += n
	}
	return req, nil
}
