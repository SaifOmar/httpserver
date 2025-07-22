package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (int, bool, error) {
	fmt.Printf("Raw incoming data to parse: %q\n", data)
	conter := 0
	read := 0
	for {
		// so now If we hit the register nurse we know we ended the headers
		if bytes.HasPrefix(data, []byte(crlf)) {
			return len(crlf), true, nil
		}
		conter++
		// fmt.Printf("conter: %d\n", conter)
		n := len(data)
		if n == 0 {
			break
		}
		crlfIndex := bytes.Index(data, []byte(crlf))
		// fmt.Printf("crlfIndex: %d\n", crlfIndex)
		if crlfIndex == -1 {
			break
		}
		sepIndex := bytes.Index(data, []byte(":"))
		if sepIndex == -1 {
			break
		}
		workingData := data[:crlfIndex]
		read += len(workingData)
		// fmt.Printf("working on data: %s\n", workingData)

		trimed := strings.TrimLeft(string(workingData[:sepIndex]), " ")
		// trimed = strings.TrimLeft(string(workingData[:sepIndex]), "\r\n")
		field_name := strings.ToLower(trimed)
		// fmt.Printf("field_name: %s\n", field_name)
		if field_name[len(field_name)-1] == ' ' {
			return 0, false, fmt.Errorf("Invalid header: %s", field_name)
		}

		if !isValidFieldValue(field_name) {
			return 0, false, fmt.Errorf("Invalid char in:%s", field_name)
		}

		value := strings.TrimSpace(string(workingData[sepIndex+1:]))

		val, exists := h[field_name]
		if exists {
			h[field_name] = val + ", " + value
		} else {
			h[field_name] = value
		}

		data = data[crlfIndex+len(crlf):]

		if len(data) <= len(crlf) {
			if string(data) == crlf {
				fmt.Printf("we came here son")
				return read + len(crlf), true, nil
			}
		}
		// fmt.Printf("new data: %s\n", data)
	}
	return read, false, nil
}
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
