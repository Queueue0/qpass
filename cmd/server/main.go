package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
)

func main() {
	cert, err := tls.LoadX509KeyPair("tls/cert.pem", "tls/key.pem")
	if err != nil {
		log.Fatal(err)
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	srv, err := tls.Listen("tcp", "localhost:1717", &config)
	if err != nil {
		log.Fatal(err)
	}

	defer srv.Close()
	for {
		c, err := srv.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go respond(c)
	}
}

func respond(c net.Conn) {
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
