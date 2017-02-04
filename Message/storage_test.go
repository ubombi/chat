package message

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestInsert(t *testing.T) {
	uuids := make([]uuid.UUID, 10)
	for i := range uuids {
		uuids[i] = uuid.NewV1()
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	storage := NewDefaultMsgStore()
	storage.Listen()
	for _, i := range rand.Perm(len(uuids)) {
		msg := NewMessage(uuid.NewV4(), fmt.Sprint(i))
		msg.ID = uuids[i]
		storage.Receive(msg)
	}
	time.Sleep(time.Second)
	for i, uid := range uuids {
		if storage.order[i] != uid {
			t.Errorf("Wrong chronological order. %v !+ %v", storage.order[0], uuids[0])
		}
	}
}
func TestLatest(t *testing.T) {
	author := uuid.NewV4()
	msgs := make([]*Message, 0, 10)
	storage := NewDefaultMsgStore()
	storage.Listen()
	if l := len(storage.Last(10)); l != 0 {
		t.Errorf("Empty storage shoud not return anything. %v results received", l)
	}
	for i := 0; i < 10; i++ {
		msg := NewMessage(author, fmt.Sprint(i))
		storage.Receive(msg)
		msgs = append(msgs, msg)
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	time.Sleep(time.Second)
	last := storage.Last(5)
	if l := len(last); l != 5 {
		t.Errorf("Wronc quantity of Messages received. 5!=%v", l)
	}

}
