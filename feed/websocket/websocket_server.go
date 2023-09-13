package websocket

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type IServer interface {
	HandleUpgrade(req *http.Request) http.Response
	HandleConnection(ws *websocket.Conn)
}

type WSClient struct {
	id     uuid.UUID
	secret string
	wsConn *websocket.Conn
}

type Server struct {
	mux     sync.Mutex
	clients []WSClient
}

func NewServer() *Server {
	return &Server{
		mux:     sync.Mutex{},
		clients: make([]WSClient, 0),
	}
}

func (s *Server) HandleUpgrade(req http.Request) *http.Response {
	defer req.Body.Close()
	id := uuid.New()

	var clientSecret []byte
	_, err := req.Body.Read(clientSecret)

	if err != nil {
		fmt.Printf("\nerror reading request body: %v\n", err)
		return nil
	}

	concatSecret := []byte(id.String() + string(clientSecret))
	sha := sha1.New().Sum(concatSecret)
	serverSecret := base64.StdEncoding.EncodeToString(sha)

	awaitingClient := WSClient{
		id:     id,
		secret: serverSecret,
		wsConn: nil,
	}

	s.clients = append(s.clients, awaitingClient)

	return &http.Response{
		StatusCode: http.StatusSwitchingProtocols,
		Status:     "101 Switching Protocols",
		Header: map[string][]string{
			"Upgrade":           {"websocket"},
			"Connection":        {"Upgrade"},
			"Sec-WebSocket-Key": {serverSecret},
		},
		Body: io.NopCloser(bytes.NewBufferString(id.String())),
	}
}

func (s *Server) HandleConnection(ws *websocket.Conn) {
	byteBuff := make([]byte, 16)
	_, err := ws.Read(byteBuff)
	if err != nil {
		ws.Close()
	}

	id, err := uuid.Parse(string(byteBuff))
	if err != nil {
		ws.Close()
	}

	var wsClient WSClient
	var wsCliIndex int
	for i, c := range s.clients {
		if c.id == id {
			s.clients[i].wsConn = ws
			wsClient = s.clients[i]
			wsCliIndex = i
			break
		}
	}

	go s.readLoop(wsClient, wsCliIndex)
}

func (s *Server) readLoop(ws WSClient, i int) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.wsConn.Read(buf)

		if err == io.EOF {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
			break
		}

		if err != nil {
			fmt.Println("Read error: ", err)
			continue
		}

		msg := buf[:n]
		fmt.Printf("Message: %s\n", msg)
		ws.wsConn.Write([]byte("Thank you for the message!"))
	}
}

func (s *Server) broadcast(b []byte) {
	for _, ws := range s.clients {
		go func(con *websocket.Conn) {
			if _, err := con.Write(b); err != nil {
				fmt.Printf("Write error: Receiver: %s \t Message: %e", con.RemoteAddr(), err)
			}
		}(ws.wsConn)
	}
}

func (s *Server) handleWSFeed(ws *websocket.Conn) {
	for {
		payload := fmt.Sprintf("data -> ")
		ws.Write([]byte(payload))
	}
}
