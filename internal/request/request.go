package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type State string

const (
	StateDone        State = "DONE"
	StateInitialized State = "INITIALIZED"
)

type Request struct {
	RequestLine RequestLine
	state       State
}

type RequestLine struct {
	Version string
	Target  string
	Method  string
}

const crlf = "\r\n"
const bufferSize = 1024

func RequestFromReader(r io.Reader) (*Request, error) {

	req := &Request{
		state: StateInitialized,
	}

	buf := make([]byte, bufferSize)
	bufLen := 0
	for req.state != StateDone {
		// grow buffer if full
		if bufLen == len(buf) {
			nb := make([]byte, len(buf)*2)
			copy(nb, buf)
			buf = nb
		}
		readN, err := r.Read(buf[bufLen:])
		if err != nil {
			return nil, fmt.Errorf("Error reading req: %w", err)
		}
		bufLen += readN
		parsedN, err := req.parse(buf[:bufLen])

		if errors.Is(err, io.EOF) {
			req.state = StateDone
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error parsing req: %w", err)
		}
		copy(buf, buf[parsedN:bufLen])
		bufLen -= parsedN
	}
	return req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	parsedN := 0
	switch r.state {
	case StateDone:
		return parsedN, errors.New("trying to read data in a done state")
	case StateInitialized:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		parsedN = n
		if n == 0 {
			break
		}
		r.RequestLine = *rl
		r.state = StateDone
	default:
		return 0, errors.New("unknown state")
	}
	return parsedN, nil
}

func parseRequestLine(stream []byte) (*RequestLine, int, error) {
	parts := strings.Split(string(stream), crlf)
	if len(parts) == 1 {
		return nil, 0, nil
	}

	firstLine := parts[0]
	lparts := strings.Split(firstLine, " ")

	if len(lparts) != 3 {
		return nil, 0, errors.New("you fucked up")
	}

	method := lparts[0]
	target := lparts[1]
	vparts := strings.Split(lparts[2], "/")
	if len(vparts) != 2 {
		return nil, 0, errors.New("invalid HTTP version")
	}
	version := vparts[1]

	n := len(firstLine)

	rl := &RequestLine{
		Target:  target,
		Version: version,
		Method:  method,
	}

	return rl, n, nil
}
