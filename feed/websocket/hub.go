package websocket

import (
	"sync"

	"github.com/milosrs/go-hls-server/feed/common"
)

type IHub interface {
	Start()
	Stop()
	Broadcast(common.Message)
	Register(crm clientRegMsg)
	Unregister(crm clientRegMsg)
}

type clientRegMsg struct {
	topic  string
	client *ClientImpl
}

type broadcastMsg struct {
	topic string
	msg   common.Message
}

type Hub struct {
	mux        sync.Mutex
	clients    map[string]map[*ClientImpl]bool
	broadcast  chan common.Message
	register   chan clientRegMsg
	unregister chan clientRegMsg
	stop       chan struct{}
}

func NewHub() IHub {
	return &Hub{
		mux:        sync.Mutex{},
		clients:    make(map[string]map[*ClientImpl]bool),
		broadcast:  make(chan common.Message, 0),
		register:   make(chan clientRegMsg, 0),
		unregister: make(chan clientRegMsg, 0),
		stop:       make(chan struct{}, 0),
	}
}

func (h *Hub) registerClient(c *clientRegMsg) {
	defer h.mux.Unlock()
	h.mux.Lock()

	if m, exists := h.clients[c.topic]; !exists {
		h.clients[c.topic] = map[*ClientImpl]bool{
			c.client: true,
		}
	} else {
		m[c.client] = true
	}
}

func (h *Hub) unregisterClient(c *clientRegMsg) {
	defer h.mux.Unlock()
	h.mux.Lock()

	if m, exists := h.clients[c.topic]; exists {
		delete(m, c.client)
		h.clients[c.topic] = m
	}
}

func (h *Hub) doBroadcast(m common.Message) {
	for client := range h.clients[m.Topic] {
		select {
		case client.send <- m:
		default:
			close(client.send)
			delete(h.clients[m.Topic], client)
		}
	}
}

func (h *Hub) Start() {
	for {
		select {
		case clientData := <-h.register:
			h.registerClient(&clientData)
		case clientData := <-h.unregister:
			h.unregisterClient(&clientData)
		case message := <-h.broadcast:
			h.doBroadcast(message)
		case <-h.stop:
			return
		}
	}
}

func (h *Hub) Stop() {
	h.stop <- struct{}{}
}

func (h *Hub) Register(crm clientRegMsg) {
	h.register <- crm
}

func (h *Hub) Unregister(crm clientRegMsg) {
	h.unregister <- crm
}

func (h *Hub) Broadcast(data common.Message) {
	h.broadcast <- data
}
