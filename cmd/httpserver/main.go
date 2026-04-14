package main

import (
	"http_protocol/internal/request"
	"http_protocol/internal/response"
	"http_protocol/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("\nServer gracefully stopped")
}

func handler(w io.Writer, r *request.Request) *server.HandlerError {
	requestTarget := r.RequestLine.RequestTarget
	switch requestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your problem is not my problem\n",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusInternalServerError,
			Message:    "Woopsie, my bad\n",
		}
	default:
		_, err := w.Write([]byte("All good, frfr\n"))
		if err != nil {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Failed to write response\n",
			}
		}
	}
	return nil
}
