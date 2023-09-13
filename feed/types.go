package feed

import (
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Message struct {
	Topic   string    `json:"topic"`
	Content []byte    `json:"content"`
	ID      uuid.UUID `json:"id"`
}

type SubChan struct {
	ID    uuid.UUID
	Chann chan Message
}

type IPubSub interface {
	Subscribe(id uuid.UUID, topic string) *SubChan
	Unsubscribe(sc SubChan, topic string)
	Publish(msg Message)
}

type IWSServer interface {
	HandleUpgrade(req http.Request) *http.Response
	HandleConnection(ws *websocket.Conn)
}
