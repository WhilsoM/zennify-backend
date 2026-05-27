package http

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) connectWs(w http.ResponseWriter, r *http.Request) {
	userID, username := userClaimsFromContext(r.Context())
	if userID == "" {
		return
	}
	if username == "" {
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	conn.SetReadLimit(512)
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if err := conn.WriteMessage(messageType, data); err != nil {
			break
		}
	}
}
