package handler

import (
	"log"
	"net/http"
	"pouic/internal/middleware"
	"pouic/internal/websocket/client"
	"github.com/gorilla/websocket"

)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func ServeWs(hub client.Hub, parser client.Parser, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	playerId, ok := middleware.PlayerIdFromContext(r.Context())
	if ok {
		client := client.NewClient(hub, parser, playerId, conn)
		client.Go()
	}
}
