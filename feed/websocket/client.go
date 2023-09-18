package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/milosrs/go-hls-server/feed/common"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client interface {
	Send(common.Message)
	SetCallback(common.OnMsgRecieved)
	Start(common.OnMsgRecieved)
	Stop()
}

type ClientImpl struct {
	hub       IHub
	conn      *websocket.Conn
	send      chan common.Message
	id        uuid.UUID
	onRecieve common.OnMsgRecieved
	topic     string
}

func NewClient(conn *websocket.Conn, hub IHub, topic string) Client {
	client := &ClientImpl{
		hub:  hub,
		conn: conn,
		send: make(chan common.Message, 0),
		onRecieve: func(msg common.Message) error {
			log.Println("UNIMPLEMENTED: ", msg)
			return nil
		},
		topic: topic,
	}

	client.hub.Register(
		clientRegMsg{
			topic:  topic,
			client: client,
		},
	)

	go client.writePump(topic)
	go client.readPump(topic)

	return client
}

// serveWs handles websocket requests from the peer.
func ServeWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}

	return conn
}

func (c *ClientImpl) SetCallback(cb common.OnMsgRecieved) {
	c.onRecieve = cb
}

func (c *ClientImpl) Send(msg common.Message) {
	c.send <- msg
}

func (c *ClientImpl) Start(cb common.OnMsgRecieved) {
	c.hub.Register(clientRegMsg{
		topic:  c.topic,
		client: c,
	})

	c.onRecieve = cb
}

func (c *ClientImpl) Stop() {
	c.hub.Unregister(clientRegMsg{
		topic:  c.topic,
		client: c,
	})
}

// Handles errors after message reading.
// Returns if read loop should be stoppped.
func (c *ClientImpl) handleError(err error) bool {
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error while reading message from client %s: %v", c.id, err)
		} else {
			log.Printf("unexpected message from client %s: %v", c.id, err)
		}
	}

	return err != nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ClientImpl) readPump(topic string) {
	defer func() {
		c.hub.Unregister(
			clientRegMsg{
				topic:  topic,
				client: c,
			},
		)
		c.conn.Close()
	}()

	pongHandler := func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	}

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(pongHandler)

	for {
		msg := common.Message{}
		err := c.conn.ReadJSON(&msg)
		if closeLoop := c.handleError(err); closeLoop {
			break
		}

		c.onRecieve(msg)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ClientImpl) writePump(topic string) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			marshalledMsg, err := json.Marshal(message)
			if err != nil {
				log.Printf("error unmarhalling message: %v", err)
			}
			w.Write(marshalledMsg)
			w.Close()

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
