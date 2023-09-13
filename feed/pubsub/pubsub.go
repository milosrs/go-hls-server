package pubsub

import (
	"sync"

	"github.com/google/uuid"
	"github.com/milosrs/go-hls-server/feed"
)

type IPubSub interface {
	Subscribe(id uuid.UUID, topic string) *feed.SubChan
	Unsubscribe(sc feed.SubChan, topic string)
	Publish(msg feed.Message)
}

type PubSub struct {
	mux    sync.Mutex
	topics map[string][]*feed.SubChan
}

func NewPubSub() IPubSub {
	return &PubSub{
		mux:    sync.Mutex{},
		topics: make(map[string][]*feed.SubChan, 0),
	}
}

func (s *PubSub) Subscribe(id uuid.UUID, topic string) *feed.SubChan {
	s.mux.Lock()
	c, ok := s.topics[topic]

	if !ok {
		c = make([]*feed.SubChan, 0)
		s.topics[topic] = c
	}
	s.mux.Unlock()

	ch := make(chan feed.Message)
	sc := &feed.SubChan{
		ID:    id,
		Chann: ch,
	}

	s.mux.Lock()
	s.topics[topic] = append(s.topics[topic], sc)
	s.mux.Unlock()

	return sc
}

func (s *PubSub) Unsubscribe(sc feed.SubChan, topic string) {
	c, ok := s.topics[topic]
	if !ok {
		return
	}

	for i, t := range c {
		if sc.ID == t.ID {
			s.topics[topic] = append(s.topics[topic][:i], s.topics[topic][i+1:]...)
			close(sc.Chann)
			break
		}
	}
}

func (s *PubSub) Publish(msg feed.Message) {
	c, ok := s.topics[msg.Topic]
	if !ok {
		return
	}

	go func() {
		for _, sc := range c {
			sc.Chann <- msg
		}
	}()
}
