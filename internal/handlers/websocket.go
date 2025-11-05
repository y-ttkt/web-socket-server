package handlers

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	connections map[*websocket.Conn]bool
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
}

var hub = Hub{
	connections: make(map[*websocket.Conn]bool),
	register:    make(chan *websocket.Conn),
	unregister:  make(chan *websocket.Conn),
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// ハンドシェイク
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()
	hub.register <- conn
	defer func() {
		hub.unregister <- conn
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		broadcastMessage(&hub, msg)
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			log.Println("Registered new connection")
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				log.Println("Removed connection")
				c.Close()
			}
		}
	}
}

// broadcastMessage 複数のクライアントに同一メッセージ返却
func broadcastMessage(h *Hub, message []byte) {
	for c := range h.connections {
		if err := c.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			c.Close()
			delete(h.connections, c)
		}
	}
}
