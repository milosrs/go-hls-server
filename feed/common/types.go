package common

import (
	"github.com/google/uuid"
)

type OnMsgRecieved func(msg Message) error

type Message struct {
	Topic   string     `json:"topic"`
	Content []byte     `json:"content"`
	ID      *uuid.UUID `json:"-"`
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
