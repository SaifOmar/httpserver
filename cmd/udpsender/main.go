package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	upd, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, upd)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("UDP connection: ", conn.RemoteAddr())
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Fatal(err)
		}
	}
}
