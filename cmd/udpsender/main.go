package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddress := "localhost:42069"

	udpAddress, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		fmt.Println("Unable to listen to connection:", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Println("Unable to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+C to exit.\n", serverAddress)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			os.Exit(1)
		}
		fmt.Print("Received: ", message)

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing message:", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", message)
	}
}
