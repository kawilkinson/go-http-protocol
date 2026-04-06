package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	linesChan := getLinesChannel(file)
	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		currLine := ""
		for {
			var bytes [8]byte
			n, err := file.Read(bytes[:])
			if err != nil {
				if currLine != "" {
					lines <- currLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("Error reading file:", err)
				break
			}

			currBytes := string(bytes[:n])
			parts := strings.Split(currBytes, "\n")
			for i, part := range parts {
				if i < len(parts)-1 {
					lines <- currLine + part
					currLine = ""
				}
			}
			currLine += parts[len(parts)-1]
		}
	}()
	return lines
}
