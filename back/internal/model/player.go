package model

type Player struct {
	Id string   `json:"id"`
	Name string `json:"name"`
}

type PlayerFull struct {
	Player
	Password string `json:"password"`
}