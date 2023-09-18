package files

import (
	"encoding/json"
	"errors"
	"log"

	gorilla_ws "github.com/gorilla/websocket"
	"github.com/milosrs/go-hls-server/feed/common"
	"github.com/milosrs/go-hls-server/feed/websocket"
	"github.com/milosrs/go-hls-server/files/model"
)

var (
	errUnknownMessage = errors.New("couldn't unmarshal data: unknown message: ")
)

const topic = "file-upload"

type FileFeed struct {
	client  websocket.Client
	hub     websocket.Hub
	service Service
}

func (f *FileFeed) Stop() {
	f.client.Stop()
}

func (f *FileFeed) onMsgRecieved(msg common.Message) error {
	var initialFileData model.InitialFileData
	var chunk model.FileChunk

	initialErr := json.Unmarshal(msg.Content, &initialFileData)
	if initialErr != nil {
		chunkErr := json.Unmarshal(msg.Content, &chunk)
		if chunkErr != nil {
			log.Printf("%v: %v", errUnknownMessage, chunkErr)
			return errUnknownMessage
		}

		_, err := f.service.AppendChunk(&chunk)
		return err
	}

	_, err := f.service.CreateFile(&initialFileData)
	return err
}

func NewFeed(ws *gorilla_ws.Conn, hub websocket.IHub, service Service) *FileFeed {
	client := websocket.NewClient(ws, hub, topic)

	ff := FileFeed{
		client:  client,
		service: service,
	}

	ff.client.Start(ff.onMsgRecieved)

	return &ff
}
