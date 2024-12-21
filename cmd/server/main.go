package main

import (
	"log"
	"net"

	"github.com/Queueue0/qpass/internal/protocol"
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
	p := protocol.Read(c)

	if p.Type() == protocol.PING {
		log.Println("PING")
		protocol.Write(c, protocol.NewPong())
	}
}
