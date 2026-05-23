package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

const MaxRetries = 5

var CRLF = []byte("\r\n")
var TERMINATOR = []byte("\r\n\r\n")

func checkAndLogError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func ParseHeaders(Req *Request, Data []byte) error {
	methods := []string{"POST", "GET"}
	parsedRequestLine := false
	var line []byte
	if Req.Headers == nil {
		Req.Headers = make(map[string]string, 0)
	}
	for {
		LineEndIndex := bytes.Index(Data, CRLF)
		if LineEndIndex == -1 {
			fmt.Println("broke")
			break
		}
		line = Data[:LineEndIndex]

		parts := bytes.Fields(line)
		if !parsedRequestLine {
			if Req.Method == "" {
				for _, method := range methods {
					if method == string(parts[0]) {
						Req.Method = string(parts[0])
					}
				}
			}
			if Req.Path == "" {
				Req.Path = string(parts[1])
			}
			if Req.Version == "" {
				Req.Version = string(parts[2])
			}
			parsedRequestLine = true
			Data = Data[LineEndIndex+len(CRLF):]
			continue
		}
		colon := bytes.IndexByte(line, ':')

		if colon == -1 {
			continue
		}

		key := bytes.TrimSpace(line[:colon])
		value := bytes.TrimSpace(line[colon+1:])

		Req.Headers[string(key)] = string(value)

		// fmt.Println("key is ", string(key))
		// fmt.Println("value is  ", string(value))

		Data = Data[LineEndIndex+len(CRLF):]
	}

	fmt.Println("path :", Req.Path)
	fmt.Println("verion:", Req.Version)
	fmt.Println("method:", Req.Method)
	return nil
}
func PraseRequest(Req *Request, Data []byte) error {
	return nil
}

func main() {
	port := "42069"
	listener, err := net.Listen("tcp", ":"+port)
	fmt.Println("Server is running on port : " + port)
	checkAndLogError(err)
	requests := []*Request{}
	for {
		conn, err := listener.Accept()
		checkAndLogError(err)
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		checkAndLogError(err)
		Req := &Request{}
		Data := []byte{}
		requests = append(requests, Req)
		BytesRead := 0
		// TODO(saif) : find \r\n\r\n
		// var lines []byte
		// FOUND_CRLF := false
		for {
			fmt.Println("hello")
			chunk := make([]byte, 1024)
			n, err := conn.Read(chunk)

			fmt.Println(n)

			if n > 0 {
				Data = append(Data, chunk[:n]...)
				BytesRead += n
				fmt.Println(string(Data))
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				checkAndLogError(err)
				break
			}
			if bytes.Contains(Data, TERMINATOR) {
				in := bytes.Index(Data, TERMINATOR)
				err = ParseHeaders(Req, Data[:in+len(CRLF)])
				for k, v := range Req.Headers {
					fmt.Printf("key : %v\n", k)
					fmt.Printf("value : %v\n", v)
					// fmt.Printf("k %v, v %s\n", k, v)
				}
				checkAndLogError(err)
				break
			}
		}
		err = PraseRequest(Req, Data)
		checkAndLogError(err)
		fmt.Println(string(Req.Body))
		for _, v := range requests {
			fmt.Println("method: ", v.Method)
			fmt.Println("path: ", v.Path)
			fmt.Println("verion: ", v.Version)
			fmt.Println("headers: ", v.Headers)
			for k, m := range v.Headers {
				fmt.Printf("key : %v\n", k)
				fmt.Printf("value : %v\n", m)
			}
			fmt.Println("body: ", string(v.Body))
		}
		conn.Close()
	}
}
