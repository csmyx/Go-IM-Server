package main

import (
	"log"

	"github.com/csmyx/Go-IM-Server/server"
)

func main() {
	log.SetPrefix("[IM] ")
	log.Println("")
	server := server.NewServer("127.0.0.1", 12345)
	server.Serve()
}
