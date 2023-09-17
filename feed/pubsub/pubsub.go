package pubsub

import (
	"sync"

	"github.com/google/uuid"
	"github.com/milosrs/go-hls-server/feed/common"
)

type PubSub struct {
	mux    sync.Mutex
	topics map[string][]*common.SubChan
}

func NewPubSub() common.IPubSub {
	return &PubSub{
		mux:    sync.Mutex{},
		topics: make(map[string][]*common.SubChan, 0),
	}
}

func (s *PubSub) Subscribe(id uuid.UUID, topic string) *common.SubChan {
	s.mux.Lock()
	c, ok := s.topics[topic]

	if !ok {
		c = make([]*common.SubChan, 0)
		s.topics[topic] = c
	}
	s.mux.Unlock()

	ch := make(chan common.Message)
	sc := &common.SubChan{
		ID:    id,
		Chann: ch,
	}

	s.mux.Lock()
	s.topics[topic] = append(s.topics[topic], sc)
	s.mux.Unlock()

	return sc
}

func (s *PubSub) Unsubscribe(sc common.SubChan, topic string) {
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

func (s *PubSub) Publish(msg common.Message) {
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
