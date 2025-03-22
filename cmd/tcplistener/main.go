package main

import (
	"fmt"
	"net"

	"github.com/agustin-carnevale/tcp-to-http/internal/request"
)

func main() {
	port := "42069"
	addr := ":" + port

	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error creating TCP listener: ", err)
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Error while accepting next tcp connection:", err)
		}
		// fmt.Println("New TCP connection:")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error parsing request", err)
		}
		printRequestLine(req.RequestLine)
		// fmt.Println("Connection closed.")
	}
}

func printRequestLine(requestLine request.RequestLine) {
	fmt.Println("Request line:")
	fmt.Println("- Method:", requestLine.Method)
	fmt.Println("- Target:", requestLine.RequestTarget)
	fmt.Println("- Version:", requestLine.HttpVersion)
}
