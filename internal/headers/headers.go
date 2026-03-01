package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	out := make(Headers)
	return out
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Look for CRLF to determine if we have a full line.
	for i := 0; i+1 < len(data); i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			// If CRLF is at the very start, headers are done.
			if i == 0 {
				return 2, true, nil
			}

			line := string(data[:i])

			// There must be a colon, and no spaces before it.
			colon := -1
			for j := 0; j < len(line); j++ {
				if line[j] == ':' {
					colon = j
					break
				}
			}
			if colon == -1 {
				return 0, false, ErrInvalidHeaderFormat()
			}
			if colon > 0 && isWhitespace(line[colon-1]) {
				return 0, false, ErrInvalidHeaderFormat()
			}

			keyRaw := trimSpace(line[:colon])
			val := ""
			if colon+1 < len(line) {
				val = trimSpace(line[colon+1:])
			}

			if keyRaw == "" {
				return 0, false, ErrInvalidHeaderFormat()
			}
			if !isValidHeaderName(keyRaw) {
				return 0, false, ErrInvalidHeaderFormat()
			}

			key := strings.ToLower(keyRaw)
			if existing, ok := h[key]; ok {
				if existing == "" {
					h[key] = val
				} else {
					h[key] = existing + ", " + val
				}
			} else {
				h[key] = val
			}
			return i + 2, false, nil
		}
	}

	// No CRLF found: not enough data yet.
	return 0, false, nil
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func ErrInvalidHeaderFormat() error {
	return errors.New("invalid header format")
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isValidHeaderName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i := 0; i < len(name); i++ {
		if !isValidTokenChar(name[i]) {
			return false
		}
	}
	return true
}

func isValidTokenChar(b byte) bool {
	if b >= 'A' && b <= 'Z' {
		return true
	}
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= '0' && b <= '9' {
		return true
	}
	switch b {
	case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	default:
		return false
	}
}
