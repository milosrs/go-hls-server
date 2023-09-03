package feed

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type Server struct {
	mux   sync.Mutex
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		mux:   sync.Mutex{},
	}
}

func (s *Server) handleConnection(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client: ", ws.RemoteAddr())

	s.mux.Lock()
	s.conns[ws] = true
	s.mux.Unlock()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)

		if err == io.EOF {
			s.mux.Lock()
			delete(s.conns, ws)
			s.mux.Unlock()
			break
		}

		if err != nil {
			fmt.Println("Read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Printf("Message: %s\n", msg)
		ws.Write([]byte("Thank you for the message!"))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(con *websocket.Conn) {
			if _, err := con.Write(b); err != nil {
				fmt.Printf("Write error: Receiver: %s \t Message: %e", con.RemoteAddr(), err)
			}
		}(ws)
	}
}

func (s *Server) handleWSFeed(ws *websocket.Conn) {
	fmt.Println("new incoming connection for feed: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("data -> ")
		ws.Write([]byte(payload))
	}
}
