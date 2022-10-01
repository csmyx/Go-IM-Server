package server

import "net"

type user struct {
	name string      // same to add by default
	addr string      // ip address
	ch   chan string // receive message sent from server to this user
	conn net.Conn
}

func newUser(conn net.Conn) *user {
	addr := conn.RemoteAddr().String()
	u := user{
		name: addr,
		addr: addr,
		ch:   make(chan string),
		conn: conn,
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
