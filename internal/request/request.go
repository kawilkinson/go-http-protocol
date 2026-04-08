package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(request)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(request []byte) (*RequestLine, error) {
	requestStr := string(request)
	requestLine := strings.Split(requestStr, "\r\n")[0]
	requestLineParts := strings.Split(requestLine, " ")
	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request line: expected 3 parts")
	}

	httpMethod := requestLineParts[0]
	if strings.ToUpper(httpMethod) != httpMethod {
		return nil, errors.New("invalid request line: method must be uppercase")
	}

	httpVersionParts := strings.Split(requestLineParts[2], "/")
	if len(httpVersionParts) != 2 {
		return nil, errors.New("invalid request line: invalid HTTP version format")
	}

	httpVersion := httpVersionParts[1]
	if httpVersion != "1.1" {
		return nil, errors.New("invalid request line: unsupported HTTP version")
	}

	return &RequestLine{
		Method:        httpMethod,
		RequestTarget: requestLineParts[1],
		HttpVersion:   httpVersion,
	}, nil
}
