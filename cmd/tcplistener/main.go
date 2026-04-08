package main

import (
	"fmt"
	"http_protocol/internal/request"
	"net"
)

func main() {
	// file, err := os.Open("messages.txt")
	// if err != nil {
	// 	fmt.Println("Error opening file:", err)
	// 	return
	// }
	// defer file.Close()

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 42069...")



	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		defer conn.Close()
		fmt.Println("Client connected:", conn.RemoteAddr())
		
		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error parsing request:", err)
			return
		}

		httpVersion := request.RequestLine.HttpVersion
		method := request.RequestLine.Method
		requestTarget := request.RequestLine.RequestTarget
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", method, requestTarget, httpVersion)

		// linesChan := getLinesChannel(conn)
		// for line := range linesChan {
		// 	fmt.Printf("%s\n", line)
		// }
	}
}

// func getLinesChannel(conn net.Conn) <-chan string {
// 	lines := make(chan string)
// 	go func() {
// 		defer close(lines)
// 		currLine := ""
// 		for {
// 			var bytes [8]byte
// 			n, err := conn.Read(bytes[:])
// 			if err != nil {
// 				if currLine != "" {
// 					lines <- currLine
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Println("Error reading file:", err)
// 				break
// 			}

// 			currBytes := string(bytes[:n])
// 			parts := strings.Split(currBytes, "\n")
// 			for i, part := range parts {
// 				if i < len(parts)-1 {
// 					lines <- currLine + part
// 					currLine = ""
// 				}
// 			}
// 			currLine += parts[len(parts)-1]
// 		}
// 	}()
// 	return lines
// }
