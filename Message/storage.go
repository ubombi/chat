package message

import (
	"encoding/binary"
	"log"
	"sort"
	"sync"

	uuid "github.com/satori/go.uuid"
)

/*
type MessageStore interface {
	Select(uuid.UUID) Message
	Slice(int, int) []*Message

	Input() <-chan Message
}
*/

type MessageReceiver interface {
	//chan<- Message
	Receive(*Message)
}

// Sorted slice, like Cassandra CK in row. With O(log(n))
type DefaultMsgStore struct {
	sync.RWMutex
	filename string
	data     map[uuid.UUID]Message
	order    []uuid.UUID
	input    chan *Message
	//order []uuid.UUID // TimeUUID (type1)
}

func (s *DefaultMsgStore) insert(msg Message) error {
	s.Lock()
	defer s.Unlock()

	if m, exists := s.data[msg.ID]; exists {
		// Update
		log.Println(m)
	} else {
		s.data[msg.ID] = msg
		time := NewTime(msg.ID)

		pos := sort.Search(len(s.order), func(i int) bool {
			return NewTime(s.order[i]) > time
		})
		// Insert into slice, without copy/ GC delete
		// https://github.com/golang/go/wiki/SliceTricks
		s.order = append(s.order, uuid.NewV1())
		copy(s.order[pos+1:], s.order[pos:])
		s.order[pos] = msg.ID
	}
	return nil

}

func (s *DefaultMsgStore) Receive(msg *Message) {
	s.input <- msg
}

func (s *DefaultMsgStore) Listen() {
	go func() {
		for msg := range s.input {
			s.insert(*msg)
		}
	}()
}

func (s *DefaultMsgStore) Last(n int) []Message {
	s.RLock()
	defer s.RUnlock()
	slice := make([]Message, 0, n)
	if l := len(s.order); l-n > 0 {
		n = l - n
	} else {
		n = 0
	}
	for _, uid := range s.order[n:] {
		slice = append(slice, s.data[uid])
	}
	return slice

}
func NewDefaultMsgStore() *DefaultMsgStore {
	return &DefaultMsgStore{
		data:  make(map[uuid.UUID]Message),
		order: make([]uuid.UUID, 0, 10),
		input: make(chan *Message, 5),
	}
}

type Time int64

func NewTime(uuid uuid.UUID) Time {
	if len(uuid) != 16 {
		return 0
	}
	time := int64(binary.BigEndian.Uint32(uuid[0:4]))
	time |= int64(binary.BigEndian.Uint16(uuid[4:6])) << 32
	time |= int64(binary.BigEndian.Uint16(uuid[6:8])&0xfff) << 48
	return Time(time)
}
