package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file: messages.txt")
	}
	defer file.Close()

	for line := range getLinesFromFileChannel(file) {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesFromFileChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	// Go routine to read lines
	go func() {
		defer close(ch)

		buffer := make([]byte, 8)
		line := ""

		for {
			n, err := f.Read(buffer)
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

			bufferString := string(buffer[:n]) // Read only the n bytes read
			parts := strings.Split(bufferString, "\n")

			if len(parts) > 1 {
				line += parts[0]
				ch <- line      // Send the completed line
				line = parts[1] // Start new line with remaining part
			} else {
				line += bufferString
			}

			// other option
			// for i := 0; i < len(parts)-1; i++ {
			// 	ch <- fmt.Sprintf("%s%s", line, parts[i])
			// 	line = ""
			// }
			// line += parts[len(parts)-1]
		}
	}()
	return ch
}
