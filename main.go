// This file is just for manual testing stuff

package main

import (
	"fmt"

	"github.com/agustin-carnevale/tcp-to-http/internal/headers"
)

func main() {

	headers := headers.Headers{}

	// data := []byte(" mykey  : romanEmpire  \r\n")
	data := []byte("   Host: localhost:42069\r\n\r\n")

	n, done, err := headers.Parse(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("HEADERS:")
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("")
	fmt.Println("Consumed bytes:")
	fmt.Println(n)
	fmt.Println("")
	fmt.Println("Done:")
	fmt.Println(done)

}
