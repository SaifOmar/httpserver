package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

func PraseRequest(Req *Request, Data []byte) error {
	return nil
}

func readFile(filename string) string {
	file, err := os.Open(filename)
	checkAndLogError(err)
	defer file.Close()
	bytes, err := io.ReadAll(file)
	checkAndLogError(err)
	return string(bytes)
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
		go handleConnection(conn, requests)
		// conn.Close()
	}
}
func SendResponse(conn net.Conn, status int, body string, headers ...map[string]string) {
	statusText := map[int]string{200: "OK", 404: "Not Found", 500: "Internal Server Error"}

	newHeaders := ""
	for _, header := range headers {
		for k, v := range header {
			newHeaders += fmt.Sprintf("%s: %s\r\n", k, v)
		}
	}
	if newHeaders == "" {
		newHeaders = "\r\n"
	}
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Length: %d\r\nConnection: close\r\n%s\r\n%s",
		status, statusText[status], len(body), newHeaders, body,
	)
	fmt.Println("response : ", response)
	conn.Write([]byte(response))
}

func handleConnection(conn net.Conn, requests []*Request) {

	ParseHeaders := func(Req *Request, Data []byte) error {
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

			Data = Data[LineEndIndex+len(CRLF):]
			if colon == -1 {
				continue
			}

			key := bytes.TrimSpace(line[:colon])
			value := bytes.TrimSpace(line[colon+1:])

			Req.Headers[string(key)] = string(value)
		}

		fmt.Println("path :", Req.Path)
		fmt.Println("verion:", Req.Version)
		fmt.Println("method:", Req.Method)
		for k, v := range Req.Headers {
			fmt.Println(k + ": " + v)
		}
		return nil
	}

	Req := &Request{}
	Data := []byte{}
	requests = append(requests, Req)
	BytesRead := 0
	in := -1
	headerLength := 0
	for {
		chunk := make([]byte, 1024)
		n, err := conn.Read(chunk)

		if n > 0 {
			Data = append(Data, chunk[:n]...)
			BytesRead += n
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			checkAndLogError(err)
			break
		}

		if bytes.Contains(Data, TERMINATOR) {
			in = bytes.Index(Data, TERMINATOR)
			headerLength = len(Data[:in+len(TERMINATOR)])
			fmt.Println("headerLength", headerLength+len(TERMINATOR))
			err = ParseHeaders(Req, Data[:in+len(CRLF)])
			if err != nil {
				panic("error")
			}
			checkAndLogError(err)
		}

		ContentLength := Req.Headers["Content-Length"]
		TransferEncoding := Req.Headers["Transfer-Encoding"]

		if ContentLength == "" && TransferEncoding == "" {
			fmt.Println("break beacause no content length or transfer encoding")
			break
		}

		if ContentLength != "" {
			ContentLengthInt, err := strconv.Atoi(ContentLength)
			checkAndLogError(err)
			if BytesRead-headerLength >= ContentLengthInt {
				Req.Body = append(Req.Body, Data[in+len(TERMINATOR):]...)
				fmt.Println("break beacause bytes read is greater than content length", BytesRead, headerLength, ContentLengthInt, n)
				break
			}
		}
	}
	switch Req.Path {
	case "/":
		SendResponse(conn, 200, "Hello World")
	case "/use-neovim-btw":
		fmt.Println("len : ", len([]byte("I use neovim-btw")))
		SendResponse(conn, 200, "I use neovim-btw")
	case "/home":
		data := `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title></title>
    <link href="css/style.css" rel="stylesheet">
  </head>
  <body style="background-color: red;">
    hello world from saif
  </body>
</html>
`
		SendResponse(conn, 200, data, map[string]string{"Content-Type": "text/html"})
	case "/index":
		SendResponse(conn, 200, readFile("index.html"), map[string]string{"Content-Type": "text/html"})
	case "/json":
		SendResponse(conn, 200, string(Req.Body), map[string]string{"Content-Type": " application/json"})
	}
	conn.Close()
}
