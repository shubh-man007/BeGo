package headers

import (
	"bytes"
	"errors"
	"strings"
)

const CRLF = "\r\n"  //Or as Prime would say a "Registered Nurse"

type Headers map[string]string

func (h Headers) Parse(data []byte) (int, bool, error) {
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
		h[key] = value

		consumed := idxCRLF + len(CRLF)
		bytesConsumed += consumed
		data = data[consumed:]
	}

	return bytesConsumed, false, nil
}
