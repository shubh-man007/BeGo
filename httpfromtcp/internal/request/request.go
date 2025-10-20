package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

const (
	stateInitialized = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	state       int
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
	idx := strings.Index(string(request), "\r\n")
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
		r.state = stateDone
		return consumed, nil

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

	req := &Request{state: stateInitialized}

	for req.state != stateDone {
		// In case buffer gets full:
		if readToIndex == len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				req.state = stateDone
				break
			}
			return nil, err
		}

		readToIndex += n

		//Parse current buffer:
		consumed, err := req.Parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			copy(buff, buff[consumed:readToIndex])
			readToIndex -= consumed
		}

	}
	return req, nil
}
