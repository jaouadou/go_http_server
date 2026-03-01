package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	add, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal("can't resolve addr", err)
	}
	conn, err := net.DialUDP("udp", nil, add)
	if err != nil {
		log.Fatal("can't dial addr", err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">>> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error reading line", err)
		}
		conn.Write([]byte(line))
	}
}
