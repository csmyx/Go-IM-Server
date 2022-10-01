package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	IP   string
	Port int
}

// create a im-server
func NewServer(ip string, port int) *Server {
	server := Server{
		IP:   ip,
		Port: port,
	}
	return &server
}

// start im-server
func (s *Server) Serve() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		// handle connection
		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	log.Println("正在处理请求中...")
}
