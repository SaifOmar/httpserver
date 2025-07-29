package headers

import (
	// "fmt"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test: Valid single header

func TestHeadersParse(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	fmt.Printf("n: %d, done: %t, err: %s\n", n, done, err)
	fmt.Printf("headers: %s\n", headers)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	//
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go \r\nSet-Person:  prime-loves-zig \r\nSet-Person:   tj-loves-ocaml \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.Equal(t, 30, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go, prime-loves-zig, tj-loves-ocaml \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	assert.Equal(t, 61, n)
	assert.False(t, done)
	//
	// headers = NewHeaders()
	// data = []byte("Set-Person: lane-loves-go, prime-loves-zig, tj-loves-ocaml \r\nHost: localhost:42069\r\n")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	// assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
	// assert.Equal(t, "localhost:42069", headers["host"])
	// assert.Equal(t, 80, n)
	// assert.False(t, done)
	//
	// headers = NewHeaders()
	// data = []byte("Set-Person: lane-loves-go, prime-loves-zig, tj-loves-ocaml")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	// assert.Equal(t, 0, n)
	// assert.False(t, done)
	//
	// headers = NewHeaders()
	// data = []byte("Host: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// require.NotNil(t, headers)
	// assert.Equal(t, "curl/7.81.0", headers["user-agent"])
	// assert.Equal(t, "localhost:42069", headers["host"])
	// assert.Equal(t, "*/*", headers["accept"])
	// assert.Equal(t, 57, n)
	// assert.False(t, done)
}
