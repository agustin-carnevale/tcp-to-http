package main

import (
	"fmt"
	"io"
	"net"
	"strings"
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
		fmt.Println("New TCP connection:")

		for line := range getLinesFromTCPChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection closed.")

	}
}

func getLinesFromTCPChannel(conn net.Conn) <-chan string {
	ch := make(chan string)

	// Go routine to read lines from TCP connn
	go func() {
		defer close(ch)

		buffer := make([]byte, 8)
		line := ""

		for {
			n, err := conn.Read(buffer)
			if err != nil {
				if err == io.EOF {
					//End of file
					if len(line) > 0 {
						ch <- line
					}
					break
				}
				fmt.Println("Error while reading file:", err)
			}

			bufferString := string(buffer[:n]) // Get only the n bytes that were read
			parts := strings.Split(bufferString, "\n")

			// Read all parts that are complete lines
			for i := 0; i < len(parts)-1; i++ {
				ch <- fmt.Sprintf("%s%s", line, parts[i]) // Send the completed line
				line = ""
			}
			line += parts[len(parts)-1] // Start new line with remaining (last) part
		}
	}()

	return ch
}
