package hub

import (
	"log"
	"pouic/internal/websocket/command"
)

// Hub registers and unregisters clients into data structure,
// stores all clients's wanted actions to distribute them to the engine.
// The engine can use the stored clients to send any messages to any player.
type Hub struct {
	players map[string]Client
	registers chan Client
	unregisters chan string
	commands chan command.Command
}

type Client interface {
	GetId() string
	Close()
}

func NewHub() *Hub {
	return &Hub{
		players: make(map[string]Client),
		registers: make(chan Client),
		unregisters: make(chan string),
		commands: make(chan command.Command),
	}
}

func (h *Hub) AddClient(c Client) {
	h.registers <- c
}

func (h *Hub) RemoveClient(id string) {
	h.unregisters <- id
}

func (h *Hub) AddCommand(command command.Command) {
	h.commands <- command
	log.Printf("command received: %v", command)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registers:
			h.players[client.GetId()] = client
		case toRemoveId := <-h.unregisters:
			h.players[toRemoveId].Close()
			delete(h.players, toRemoveId)
		}
	}
}