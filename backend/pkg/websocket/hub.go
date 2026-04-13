package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/wafi11/workspaces/config"
)

type Hub struct {
	clients map[string][]*websocket.Conn
	conf    *config.Config
	mu      sync.RWMutex
}

func NewHub(conf *config.Config) *Hub {
	return &Hub{
		conf:    conf,
		clients: make(map[string][]*websocket.Conn),
	}
}
func (h *Hub) Handler(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		log.Println("[ws] no token provided")
		return echo.ErrUnauthorized
	}

	data, err := config.ValidationToken(token, h.conf)
	if err != nil {
		log.Printf("[ws] invalid token: %v", err)
		return echo.ErrUnauthorized
	}

	log.Printf("[ws] client connecting: userID=%s", data.UserID)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("[ws] upgrade failed: %v", err)
		return err
	}

	h.Register(data.UserID, conn)
	log.Printf("[ws] client registered: userID=%s total_clients=%d", data.UserID, len(h.clients))
	defer func() {
		h.Unregister(data.UserID, conn)
		log.Printf("[ws] client disconnected: userID=%s", data.UserID)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[ws] read error userID=%s: %v", data.UserID, err)
			break
		}
	}
	return nil
}

func (h *Hub) Register(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], conn)
}

func (h *Hub) Unregister(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns := h.clients[userID]
	for i, c := range conns {
		if c == conn {
			h.clients[userID] = append(conns[:i], conns[i+1:]...)
			conn.Close()
			break
		}
	}
}

func (h *Hub) SendToUser(userID string, payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	msg, _ := json.Marshal(payload)
	for _, conn := range h.clients[userID] {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
