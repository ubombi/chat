package message

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type Message struct {
	ID      uuid.UUID //TimeUUID
	Text    string
	Author  uuid.UUID
	Version uuid.UUID //TimeUUID
}

func (m *Message) String() string {
	return fmt.Sprintf("FROM: %v \nTEXT:\n%v\n", m.Author, m.Text)
}

// Update text and set version TimeUUID tag
func (m *Message) Update(text string) {
	m.Version = uuid.NewV1()
	m.Text = text
}
func NewMessage(author uuid.UUID, text string) *Message {
	id := uuid.NewV1()
	return &Message{
		ID:      id,
		Text:    text,
		Author:  author,
		Version: id,
	}
}
