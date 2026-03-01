package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/http_server/request"
)

var PORT int = 42069

func main() {
	portStr := fmt.Sprintf(":%s", strconv.Itoa(PORT))
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatal("couldn't connect", err)
	}
	defer l.Close()

	fmt.Printf("listening on: localhost:%d\n", PORT)

	// lines := getLinesChannel(file)
	// line := <-lines

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("error accepting connetction: ", err)
		}
		fmt.Println("connection accpeted!")

		request, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("erro handeling request: ", err)
		}

		rl := request.RequestLine
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
			rl.Method,
			rl.Target,
			rl.Version,
		)

	}
}
