package server

import "net"

type user struct {
	name   string      // same to add by default
	addr   string      // ip address
	ch     chan string // receive message sent from server to this user
	conn   net.Conn
	server *Server
}

func newUser(conn net.Conn, s *Server) *user {
	addr := conn.RemoteAddr().String()
	u := user{
		name:   addr,
		addr:   addr,
		ch:     make(chan string),
		conn:   conn,
		server: s,
	}

	go u.listenMessage()

	return &u
}

func (u *user) listenMessage() {
	for {
		msg := <-u.ch

		u.conn.Write([]byte(msg + "\n"))
		Debug.Println(msg)
	}
}

func (u *user) String() string {
	str := "[" + u.addr + "]" + u.name
	return str
}

func (u *user) chatFmt() string {
	return u.name + ": "
}

// broadcast login message
func (u *user) login() {
	msg := u.chatFmt() + "已上线"
	u.server.broadcast(msg)
	Log.Println(u.String(), "LOGIN")
}

// broadcast logout message
func (u *user) logout() {
	msg := u.chatFmt() + "已下线"
	u.server.broadcast(msg)
	Log.Println(u.String(), "LOGOUT")
}

// broadcast chat message
func (u *user) broadcastChat(m string) {
	msg := u.chatFmt() + m
	u.server.broadcast(msg)
	Log.Println(u.String(), "BROAD CHAT")
}

// broadcast messages to all online users via channel
func (s *Server) broadcast(msg string) {
	s.msgCh <- msg
	Debug.Println("brodcast", msg)
}
