package main

import (
	"bytes"
	"fmt"
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
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Println(line)
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
