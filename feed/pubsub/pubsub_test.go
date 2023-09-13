package pubsub_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/milosrs/go-hls-server/feed"
	"github.com/milosrs/go-hls-server/feed/pubsub"
)

func TestWholeProcess(t *testing.T) {
	ps := pubsub.NewPubSub()
	uuid := uuid.New()
	topic := "test"
	message := feed.Message{
		Topic:   topic,
		Content: []byte("Hello world of testing"),
		ID:      uuid,
	}

	newSubChan := ps.Subscribe(uuid, topic)
	assert.NotEqual(t, newSubChan, nil)

	ps.Publish(message)

	data := <-newSubChan.Chann

	assert.Equal(t, data.ID, message.ID)
	assert.Equal(t, data.Content, message.Content)
	time.Sleep(1000)

	ps.Unsubscribe(*newSubChan, topic)

	data, ok := <-newSubChan.Chann
	assert.Equal(t, ok, false)
}

func TestAsyncProcess(t *testing.T) {
	const agents = 100

	ps := pubsub.NewPubSub()
	uuid := uuid.New()
	topic := "test"
	message := feed.Message{
		Topic:   topic,
		Content: []byte("Hello world of testing"),
		ID:      uuid,
	}
	subChans := make([]*feed.SubChan, 0)

	creationWG := sync.WaitGroup{}
	creationWG.Add(agents)
	for i := 0; i < agents; i++ {
		go func() {
			subChans = append(subChans, ps.Subscribe(uuid, topic))
			creationWG.Done()
		}()
	}
	creationWG.Wait()

	for i := 0; i < agents; i++ {
		go func() {
			ps.Publish(message)
		}()
	}

	wg := sync.WaitGroup{}
	wg.Add(agents)

	go func() {
		for _, sc := range subChans {
			<-sc.Chann
			wg.Done()
		}
	}()

	wg.Wait()
}
