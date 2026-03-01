package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("couldn't connect", err)
	}
	defer l.Close()

	// lines := getLinesChannel(file)
	// line := <-lines

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("error accepting connetction", err)
		}
		fmt.Println("connection accpeted!")
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		line := ""
		buf := make([]byte, 8)

		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			data := buf[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				line += string(data[:i])
				ch <- line
				line = ""
				data = data[i+1:]
			}

			line += string(data)
		}
		if line != "" {
			ch <- line
		}
	}()
	return ch
}
