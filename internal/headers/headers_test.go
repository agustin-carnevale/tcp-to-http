package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestHeaders(t *testing.T) {

	// Test: Valid single header
	headers := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = Headers{}
	data = []byte("    Host:    localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 35, n)
	assert.False(t, done)

	// Test: Valid 2 headers
	headers = Headers{}
	data = []byte("Host: localhost:42069\r\n Content-Type: application/json \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid Done
	headers = Headers{}
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.True(t, done)
	assert.Equal(t, 2, n)

	// Test: Invalid spacing header
	headers = Headers{}
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid header key with digits and special character
	headers = Headers{}
	data = []byte("Host_!19:    localhost:42069     \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host_!19"])
	assert.Equal(t, 35, n)
	assert.False(t, done)

	// Test: Invalid special characters on header key
	headers = Headers{}
	data = []byte("H©st@: localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valiud multiple same key headers
	headers = Headers{}
	headers["set-developer"] = "Agustin"
	data = []byte("Set-Developer: Michael \r\n\r\n")
	_, _, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "Agustin, Michael", headers["set-developer"])
}
