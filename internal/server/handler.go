package server

import (
	"http_protocol/internal/request"
	"http_protocol/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message string	
}

type Handler func(w io.Writer, request *request.Request) *HandlerError

func (e *HandlerError) WriteError(w io.Writer) {
	response.WriteStatusLine(w, e.StatusCode)
	errorHeaders := response.GetDefaultHeaders(len(e.Message))
	response.WriteHeaders(w, errorHeaders)
	w.Write([]byte(e.Message))
}
