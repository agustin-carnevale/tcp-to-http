package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Errorf("Error opening file: messages.txt")
	}
	defer file.Close()

	buffer := make([]byte, 8)

	for {
		_, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				//End of file
				break
			}
			fmt.Println("Error while reading file:", err)
			return
		}
		fmt.Printf("read: %s\n", string(buffer))
	}
}
