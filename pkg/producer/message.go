package producer

import (
	"github.com/satori/go.uuid"
)

type Message struct {
	ID   uuid.UUID `json:"id"`
	Body []byte    `json:"body"`
}

func NewMessage(body []byte) *Message {
	return &Message{uuid.NewV4(), body}
}

func (m Message) String() string {
	return m.ID.String()
}
