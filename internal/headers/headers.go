package headers

import (
	"bytes"
	"fmt"
	// "strings"
	"unicode"
)

type Headers map[string]string

const (
	CRLF        = "\r\n"
	HEADERS_END = "\r\n\r\n"
)

// I am just going to redesign this to not be a loop
// I will read from the reader and only parse here if I can parse (a haeder line is complete) and return the number of bytes read
// else I will return the 0
func (h Headers) Parse(data []byte) (int, bool, error) { // TODO: this is a mess

	// if bytes.HasPrefix(data, []byte(HEADERS_END)) {
	// 	return len(HEADERS_END), true, nil
	// }

	crlfIndex := bytes.Index(data, []byte(CRLF))
	if crlfIndex == -1 {
		return 0, false, nil
	}
	if crlfIndex == 0 {
		return len(CRLF), false, nil
	}
	header_line := data[:crlfIndex]
	sepIndex := bytes.Index(header_line, []byte(":"))
	if sepIndex == -1 {
		return 0, false, fmt.Errorf("Invalid header: %s", header_line)
	}
	field_name := bytes.ToLower(bytes.TrimLeft(header_line[:sepIndex], " "))
	if field_name[len(field_name)-1] == ' ' {
		return 0, false, fmt.Errorf("Invalid header: %s", field_name)
	}
	if !isValidFieldValue(string(field_name)) {
		return 0, false, fmt.Errorf("Invalid char in:%s", string(field_name))
	}
	value := bytes.TrimSpace(header_line[sepIndex+1:])
	val, exists := h[string(field_name)]

	read := len(data[:crlfIndex]) + len(CRLF)
	if exists {
		h[string(field_name)] = val + ", " + string(value)
		read = len(val) + len(value) + len(CRLF)
	} else {
		h[string(field_name)] = string(value)
	}
	return read, false, nil

}

// func (h Headers) Parse(data []byte) (int, bool, error) {
// 	read := 0
// 	for {
// 		fmt.Printf("Data: %q\n", data)
// 		if bytes.HasPrefix(data, []byte(HEADERS_END)) {
// 			return len(HEADERS_END), true, nil
// 		}
//
// 		crlfIndex := bytes.Index(data, []byte(CRLF))
// 		// this is just saying I can't parse this give me more data
// 		if crlfIndex == -1 {
// 			break
// 		}
// 		// you give me crlf at the beggeing I say yaya okay give me more without it next time
// 		if crlfIndex == 0 {
// 			return len(CRLF), false, nil
// 		}
// 		header_line := data[:crlfIndex]
// 		sepIndex := bytes.Index(header_line, []byte(":"))
// 		if sepIndex == -1 {
// 			break
// 		}
// 		field_name := bytes.ToLower(bytes.TrimLeft(header_line[:sepIndex], " "))
//
// 		if field_name[len(field_name)-1] == ' ' {
// 			return 0, false, fmt.Errorf("Invalid header: %s", field_name)
// 		}
// 		if !isValidFieldValue(string(field_name)) {
// 			return 0, false, fmt.Errorf("Invalid char in:%s", string(field_name))
// 		}
// 		value := bytes.TrimSpace(header_line[sepIndex+1:])
// 		val, exists := h[string(field_name)]
// 		if exists {
// 			h[string(field_name)] = val + ", " + string(value)
// 		} else {
// 			h[string(field_name)] = string(value)
// 		}
//
// 		read += len(data[:crlfIndex])
// 		data = data[crlfIndex+len(CRLF):]
// 		if len(data) <= len(CRLF) {
// 			if string(data) == CRLF {
// 				return read + len(CRLF), true, nil
// 			} else {
// 				return read, false, nil
// 			}
// 		}
//
// 	}
//
// 	return 0, false, nil
// }

func NewHeaders() Headers {
	return Headers{}
}

func isValidFieldValue(value string) bool {
	specialChars := map[rune]bool{
		'~':  true,
		'|':  true,
		'`':  true,
		'_':  true,
		'^':  true,
		'.':  true,
		'-':  true,
		'+':  true,
		'*':  true,
		'\'': true,
		'&':  true,
		'%':  true,
		'$':  true,
		'#':  true,
		'!':  true,
	}
	for _, c := range value {
		if !unicode.IsLower(c) && !unicode.IsDigit(c) && !specialChars[c] {
			fmt.Printf("c: %c\n", c)
			return false
		}

	}
	return true
}
