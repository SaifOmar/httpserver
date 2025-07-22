package main

import (
	"bytes"
	"fmt"
	"httpserver/internal/request"
	"io"
	"log"
	"net"
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

		fmt.Printf("Incoming connection to the listner: %v", conn)
		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")

		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		conn.Close()
		fmt.Println("Connection has been closed")

	}

}
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
