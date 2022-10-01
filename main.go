package main

import (
	"github.com/csmyx/Go-IM-Server/server"
)

func main() {
	server := server.NewServer("127.0.0.1", 12345)
	server.Serve()
}
