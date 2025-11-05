package headers

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const CRLF = "\r\n"

var fieldNameConstraint = regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.^_` + "`" + `|~]+$`)

type Headers struct {
	header map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		header: map[string]string{},
	}
}

func (h *Headers) Get(name string) string {
	return h.header[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	v, ok := h.header[strings.ToLower(name)]
	if !ok {
		h.header[strings.ToLower(name)] = value
	} else {
		h.header[strings.ToLower(name)] = fmt.Sprintf("%s,%s", v, value)
	}
}

func (h *Headers) Replace(name, value string) {
	h.header[strings.ToLower(name)] = value
}

func (h *Headers) Iter() map[string]string {
	return h.header
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	bytesConsumed := 0

	for len(data) > 0 {
		idxCRLF := bytes.Index(data, []byte(CRLF))
		if idxCRLF == -1 {
			return bytesConsumed, false, nil
		}

		if idxCRLF == 0 {
			return bytesConsumed + len(CRLF), true, nil
		}

		line := strings.TrimSpace(string(data[:idxCRLF]))
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 || (colonIdx > 0 && line[colonIdx-1] == ' ') {
			return bytesConsumed, false, errors.New("invalid field-line syntax")
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		if !fieldNameConstraint.Match([]byte(key)) {
			return bytesConsumed, false, errors.New("invalid tchar for field-name")
		}

		h.Set(key, value)

		consumed := idxCRLF + len(CRLF)
		bytesConsumed += consumed
		data = data[consumed:]
	}

	return bytesConsumed, false, nil
}
