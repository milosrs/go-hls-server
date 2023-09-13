package feed

import "github.com/google/uuid"

type Message struct {
	Topic   string    `json:"topic"`
	Content []byte    `json:"content"`
	ID      uuid.UUID `json:"id"`
}

type SubChan struct {
	ID    uuid.UUID
	Chann chan Message
}
