package message

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestSum256(t *testing.T) {
	var sum1, sum2 [32]byte

	msg := NewMessage(uuid.NewV4(), "Test message string")
	sum1 = msg.Sum256()
	msg.Signature = []byte("ShoudNotBeUsedInHashSum")
	sum2 = msg.Sum256()
	if sum1 != sum2 {
		t.Error("Signature field shoud not be used in checksum")
	}
	msg.Update("New text string")
	t.Log(msg)
	sum2 = msg.Sum256()
	if sum1 == sum2 {
		t.Logf("Hash1: %v, Hash2: %v", sum1, sum2)
		t.Error("Same checksum for different messages")
	}
}

func TestUpdate(t *testing.T) {
	msg := NewMessage(uuid.NewV4(), "Text1")
	v := msg.Version
	id := msg.ID
	msg.Update("Text1")
	if v != msg.Version {
		t.Error("Version missmatch. Message was not changed.")
	}
	msg.Update("Text2")
	if v == msg.Version {
		t.Error("Version was not updated")
	}
	if msg.Text != "Text2" {
		t.Error("Text was not updated")
	}
	if msg.ID != id {
		t.Error("Message ID shoud not change.")
	}
}
