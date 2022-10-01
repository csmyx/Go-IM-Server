package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const msgBufLen = 4096

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

	go s.listenMessage()

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

	s.login(u)

	// receive message sent by user
	go func() {
		buf := make([]byte, msgBufLen)
		for {
			n, err := u.conn.Read(buf)
			if n == 0 {
				s.logout(u)
				return
			}
			if err != nil && err != io.EOF {
				Log.Fatalln(err)
			}

			// strip last newline
			msg := string(buf[:n-1])
			s.broadcastChat(u, msg)
		}
	}()

	select {}
}

// broadcast login message
func (s *Server) login(u *user) {
	msg := u.chatFmt() + "已上线"
	s.broadcast(msg)
	Log.Println(u.String(), "LOGIN")
}

// broadcast logout message
func (s *Server) logout(u *user) {
	msg := u.chatFmt() + "已下线"
	s.broadcast(msg)
	Log.Println(u.String(), "LOGOUT")
}

// broadcast chat message
func (s *Server) broadcastChat(u *user, m string) {
	msg := u.chatFmt() + m
	s.broadcast(msg)
	Log.Println(u.String(), "BROAD CHAT")
}

// broadcast messages to all online users via channel
func (s *Server) broadcast(msg string) {
	s.msgCh <- msg
	Debug.Println("brodcast", msg)
}

// listening message channel
func (s *Server) listenMessage() {
	for {
		msg := <-s.msgCh
		s.mtx.RLock()
		for _, u := range s.onlineUsers {
			u.ch <- msg
		}
		s.mtx.RUnlock()
	}
}
