package server

import (
	"fmt"
	"http_protocol/internal/request"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	done     atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Server is listening on port %d...\n", port)
	server := &Server{listener: listener, done: atomic.Bool{}}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.done.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		if s.done.Load() {
			return
		}

		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	_, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	// static response for now
	response := "HTTP/1.1 200 OK\n" +
		"Content-Type: text/plain\n" +
		"Content-Length: 13\n" +
		"\n" +
		"Hello World! \n"
	conn.Write([]byte(response))

	fmt.Printf("%v", response)
}
