package server

import (
	"bytes"
	"fmt"
	"http_protocol/internal/request"
	"http_protocol/internal/response"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	done     atomic.Bool
	handler Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Server is listening on port %d...\n", port)
	server := &Server{listener: listener, done: atomic.Bool{}, handler: handler}
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

	request, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.WriteError(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, request)
	if hErr != nil {
		hErr.WriteError(conn)
		return
	}

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Println("Error writing status line:", err)
		return
	}
	responseHeaders := response.GetDefaultHeaders(len(buf.Bytes()))
	err = response.WriteHeaders(conn, responseHeaders)
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}

	conn.Write(buf.Bytes())
}
