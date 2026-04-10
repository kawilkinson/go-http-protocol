package headers

import (
	"errors"
	"strings"
)

const crlf = "\r\n"

var ValidHeaderKeyRunes = map[rune]struct{}{
	'!': {}, '#': {}, '$': {}, '%': {}, '&': {}, '*': {}, '+': {}, '-': {}, '.': {}, '^': {}, '_': {}, '`': {}, '|': {}, '~': {},

	'0': {}, '1': {}, '2': {}, '3': {}, '4': {}, '5': {}, '6': {}, '7': {}, '8': {}, '9': {},

	'A': {}, 'B': {}, 'C': {}, 'D': {}, 'E': {}, 'F': {}, 'G': {}, 'H': {}, 'I': {}, 'J': {}, 'K': {}, 'L': {}, 'M': {}, 'N': {}, 'O': {}, 'P': {},
	'Q': {}, 'R': {}, 'S': {}, 'T': {}, 'U': {}, 'V': {}, 'W': {}, 'X': {}, 'Y': {}, 'Z': {},

	'a': {}, 'b': {}, 'c': {}, 'd': {}, 'e': {}, 'f': {}, 'g': {}, 'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {}, 'n': {}, 'o': {}, 'p': {},
	'q': {}, 'r': {}, 's': {}, 't': {}, 'u': {}, 'v': {}, 'w': {}, 'x': {}, 'y': {}, 'z': {},
}

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
	if headerName != strings.TrimSpace(headerName) {
		return 0, false, errors.New("invalid header: header name contains spaces")
	}

	headerName = strings.TrimSpace(headerName)
	if !isValidHeaderKey(headerName) {
		return 0, false, errors.New("invalid header: header name contains invalid characters")
	}

	headerName = strings.ToLower(headerName)
	headerValue := strings.TrimSpace(dataParts[1])
	if existingValue, exists := h[headerName]; exists {
		headerValue = existingValue + ", " + headerValue
	}
	h[headerName] = headerValue
	return idx + 2, false, nil
}

func isValidHeaderKey(key string) bool {
	for _, char := range key {
		if _, exists := ValidHeaderKeyRunes[char]; !exists {
			return false
		}
	}
	return true
}
