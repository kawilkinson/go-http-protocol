package request

import (
	"errors"
	"fmt"
	"http_protocol/internal/headers"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
	for request.state != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}

		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.state, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	requestStr := string(request)
	idx := strings.Index(requestStr, crlf)
	if idx == -1 {
		return nil, 0, nil
	}
	requestStr = requestStr[:idx]
	requestLine, err := parseRequestLineFromString(requestStr[:idx])
	if err != nil {
		return nil, 0, err
	}

	return requestLine, len(requestStr) + 2, nil
}

func parseRequestLineFromString(str string) (*RequestLine, error) {
	requestLine := strings.Split(str, "\r\n")[0]

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

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// need more data to parse, break and wait for next read
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	// initialized state
	case requestStateInitialized:
		requestLine, bytesRead, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if requestLine == nil {
			return 0, nil
		}
		if bytesRead == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return bytesRead, nil

	// parsing headers state
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateDone
		}
		return n, nil

	// done state
	case requestStateDone:
		return 0, errors.New("error: trying to read data in a done state")
	default:
		return 0, errors.New("error: unknown state")
	}
}
