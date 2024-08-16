package main

import (
	"crypto/tls"
	"log"
)

func send(s string) string {
	conf := tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := tls.Dial("tcp", "localhost:1717", &conf)
	if err != nil {
		log.Fatal(err)
	}

	c.Write([]byte(s))

	buffer := make([]byte, 1024)
	_, err = c.Read(buffer)

	return string(buffer[:])
}
