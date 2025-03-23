// This file is just for manual testing stuff

package main

import (
	"fmt"
	"io"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func main() {

	// headers := headers.Headers{}

	// // data := []byte(" mykey  : romanEmpire  \r\n")
	// data := []byte("   Host: localhost:42069\r\n\r\n")

	// n, done, err := headers.Parse(data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("HEADERS:")
	// for key, value := range headers {
	// 	fmt.Printf("%s: %s\n", key, value)
	// }
	// fmt.Println("")
	// fmt.Println("Consumed bytes:")
	// fmt.Println(n)
	// fmt.Println("")
	// fmt.Println("Done:")
	// fmt.Println(done)

	// reader := &chunkReader{
	// 	data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
	// 	numBytesPerRead: 3,
	// }
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 16\r\n" +
			"\r\n" +
			"This is the body",
		numBytesPerRead: 3,
	}

	r, err := request.RequestFromReader(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	r.Print()

}
