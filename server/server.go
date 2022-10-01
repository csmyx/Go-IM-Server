package server

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP          string
	Port        int
	onlineUsers map[string]*user
	mtx         sync.RWMutex
	msgCh       chan string // send messages to online users
}

// create an im-server
func NewServer(ip string, port int) *Server {
	s := Server{
		IP:          ip,
		Port:        port,
		onlineUsers: make(map[string]*user),
		msgCh:       make(chan string),
	}
	return &s
}

// start im-server
func (s *Server) Serve() {
	// socket listen
	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		Log.Fatalln(err)
	}
	defer listener.Close()

	Log.Printf("启动服务器[%s] SUCCEED\n", addr)

	go s.broadcast()

	for {
		// accept connection
		conn, err := listener.Accept()
		if err != nil {
			Log.Println("接受连接请求 FAILED:", err)
			continue
		}

		// handle connection
		go s.handler(conn)
	}
}

func (s *Server) handler(conn net.Conn) {
	u := newUser(conn)
	s.mtx.Lock()
	s.onlineUsers[u.name] = u
	s.mtx.Unlock()

	Log.Printf("处理请求[%s] STARTED\n", u.addr)

	s.Login(u)

	select {}
}

// broadcast login message
func (s *Server) Login(u *user) {
	msg := "[" + u.addr + "]" + u.name + " 已上线"
	s.msgCh <- msg
	Log.Println(msg)
}

// broadcast messages to all online users
func (s *Server) broadcast() {
	for {
		msg := <-s.msgCh
		s.mtx.RLock()
		for _, u := range s.onlineUsers {
			u.ch <- msg
		}
		s.mtx.RUnlock()
		Debug.Println("brodcast", msg)
	}
}
