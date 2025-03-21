package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	port := "42069"
	addr := "localhost:" + port

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving UDP address: ", err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection: ", err)
		return
	}
	defer udpConn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			return
		}
		udpConn.Write([]byte(line))
		if err != nil {
			fmt.Println("Error sending UDP packet:", err)
			return
		}
	}
}
