package request

import (
	"bytes"
	"fmt"
	"httpserver/internal/headers"
	"io"
	"strings"
	"unicode"
)

type parserState string

const (
	StateInit         parserState = "init"
	StateDone         parserState = "done"
	StateParseHeaders parserState = "parsing_headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parserState
}

func NewRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// so essentially how this is working is that we are reading from the reader which pushes chuncks of data to us
// we are keeping track of the index of the buffer and the amount of data we have read so far
// we get the number of bytes read and then we add it the the index of the buffer to get the new index of in the buffer
// then we send the date to the parser which will parse the data and return either 0 if it parsed and found a rigester nurse
// or the index of where it stopped parsing because it didn't find a rigester nurse
func RequestFromReader(reader io.Reader) (*Request, error) {
	r := NewRequest()
	buff := make([]byte, 1024)
	buffIndex := 0
	for r.state != StateDone {
		readN, err := reader.Read(buff[buffIndex:])
		if err != nil {
			return nil, fmt.Errorf("Error reading from io.Reader: %s", err)
		}

		buffIndex += readN
		parseN, err := r.parse(buff[:buffIndex])
		if err != nil {
			return nil, fmt.Errorf("Error parsing: %s", err)
		}
		copy(buff, buff[parseN:buffIndex])
		buffIndex -= parseN

	}
	// note here parseN might be 0 or the index of the start to retry the parsing
	return r, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		if r.state == StateInit {
			line, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *line
			// fmt.Printf("line: %s\n", r.RequestLine)
			read += n
			r.state = StateParseHeaders
			return read, nil

		}

		if r.state == StateParseHeaders {
			fmt.Printf("Data: %q\n", data[read:])
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}
			read += n
			if done {
				r.state = StateDone
				return read, nil
			}
			return read, nil

		}
		if r.state == StateDone {
			break outer
		}
	}
	return read, nil
}

func parseRequestLine(line []byte) (*RequestLine, int, error) {
	sep := []byte("\r\n")
	// get the index of the rigester nurse
	index := bytes.Index(line, sep)
	if index == -1 {
		return nil, 0, nil
	}

	// get the start ine wich will be everything before the \r\n
	startLine := line[:index]
	// get the rest of the message which is everything after the \r\n
	read := index + (len("\r\n"))

	parts := strings.Split(string(startLine), " ")

	if len(parts) != 3 {
		return nil, read, fmt.Errorf("Invalid start line: %s", line)

	}

	method := parts[0]
	requestTarget := parts[1]
	httpVersion := strings.Split(parts[2], "/")[1]
	for i := range method {
		if !unicode.IsLetter(rune(method[i])) {
			return nil, read, fmt.Errorf("Invalid method: %s", method)
		}
		if !unicode.IsUpper(rune(method[i])) {
			return nil, read, fmt.Errorf("Invalid method: %s", method)
		}
	}
	if httpVersion != "1.1" {
		return nil, read, fmt.Errorf("Invalid HTTP version: %s", httpVersion)
	}
	// log.Printf("method: %s, requestTarget: %s, httpVersion: %s", method, requestTarget, httpVersion)
	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   httpVersion,
	}, read, nil
}
