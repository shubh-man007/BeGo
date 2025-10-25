package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldLineParser(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 25, n)
	assert.True(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host:          localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 48, n)
	assert.True(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers.Set("Existing", "header")
	data = []byte("Host: localhost:42069\r\nContent-Type: application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, "application/json", headers.Get("Content-type"))
	assert.Equal(t, "header", headers.Get("Existing"))
	assert.Equal(t, 57, n)
	assert.True(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Caps Header
	headers = NewHeaders()
	data = []byte("HOST:  LOcalhoST:    8080\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "LOcalhoST:    8080", headers.Get("Host"))
	assert.Equal(t, 27, n)
	assert.False(t, done)

	// Test: Invalid Header
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple Values
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 86, n)
	assert.True(t, done)
	assert.Equal(t, "lane-loves-go,prime-loves-zig,tj-loves-ocaml", headers.Get("Set-Person"))
}

// go test ./...
