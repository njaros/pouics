package handler

import (
	"context"
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

func ServeWs(ctx context.Context, hub client.Hub, parser client.Parser, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	playerId, ok := middleware.PlayerIdFromContext(ctx)
	if ok {
		client := client.NewClient(hub, parser, playerId, conn)
		client.Go()
	}
}
