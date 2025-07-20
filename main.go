package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Println(line)
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
