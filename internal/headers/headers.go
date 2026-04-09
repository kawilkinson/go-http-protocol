package headers

import (
	"errors"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	dataStr := string(data)
	if !strings.Contains(dataStr, crlf) {
		return 0, false, errors.New("invalid header: missing CRLF")
	}
	if dataStr[:2] == crlf {
		// empty line, headers are done, consume the CRLF
		return 2, true, nil
	}

	idx := strings.Index(dataStr, crlf)
	dataParts := strings.SplitN(dataStr[:idx], ":", 2)
	headerName := dataParts[0]
	if headerName != strings.TrimRight(headerName, " ") {
		return 0, false, errors.New("invalid header: header name contains spaces")
	}

	headerName = strings.TrimSpace(headerName)
	headerValue := strings.TrimSpace(dataParts[1])
	h[headerName] = headerValue
	return idx + 2, false, nil
}
