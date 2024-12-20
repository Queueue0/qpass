package main

import (
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
)

func main() {
	srv, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	defer srv.Close()

	for {
		c, err := srv.Accept()
		if err != nil {
			panic(err)
		}

		go respond(c)
	}
}

func respond(c net.Conn) {
	defer c.Close()
	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	i := slices.Index(buf, byte(0))
	if i == -1 {
		i = 1024
	}

	text := strings.TrimSpace(string(buf[:i]))
	resp := fmt.Sprintf("Request: \"%v\" Response: \"received\"\n", text)
	log.Println(resp)
	c.Write([]byte(resp))
	c.Close()
}
