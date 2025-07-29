package main

import (
	"bytes"
	"fmt"
	"httpserver/internal/request"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", ":42069") // prime's got taste
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Incoming connection to the listener: %v\n", conn.RemoteAddr())

		// Set timeout in case client hangs
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Failed to read request: %v\n", err)
			conn.Close()
			continue
		}

		fmt.Printf("Method: %s, RequestTarget: %s, HTTP Version: %s\n",
			request.RequestLine.Method,
			request.RequestLine.RequestTarget,
			request.RequestLine.HttpVersion,
		)
		fmt.Printf("Headers: %s\n", request.Headers)
		for k, v := range request.Headers {
			fmt.Printf("Header: %s, Value: %s\n", k, v)
		}

		conn.Close()
		fmt.Println("Connection has been closed")
	}
}

//	func main() {
//		l, err := net.Listen("tcp", ":42069") // prime's got taste
//		if err != nil {
//			log.Fatal(err)
//		}
//		defer l.Close()
//		for {
//			conn, err := l.Accept()
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Printf("Incoming connection to the listner: %v", conn)
//			request, err := request.RequestFromReader(conn)
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Printf("Method: %s, RequestTarget: %s, HTTP Version: %s\n", request.RequestLine.Method, request.RequestLine, request.RequestLine.HttpVersion)
//			conn.Close()
//			fmt.Println("Connection has been closed")
//
//		}
//
// }
func getLinesChannel(file io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer file.Close()
		defer close(out)
		str := ""
		for {
			data := make([]byte, 8)
			_, err := file.Read(data)
			if err != nil {
				log.Fatal(err)
				break
			}
			i := bytes.IndexByte(data, '\n')
			if i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""

			}
			str += string(data)

		}
		if len(str) != 0 {
			out <- str
		}
	}()
	return out
}
