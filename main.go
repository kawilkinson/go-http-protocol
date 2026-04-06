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

	currLine := ""
	for {
		var bytes [8]byte
		n, err := file.Read(bytes[:])
		if err != nil {
			if currLine != "" {
				fmt.Printf("read: %s\n", currLine)
				currLine = ""
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
				fmt.Printf("read: %s%s\n", currLine, part)
				currLine = ""
			}
		}
		currLine += parts[len(parts)-1]
	}
}
