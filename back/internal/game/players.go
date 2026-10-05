package game

import "pouic/internal/model"

type Socket int

type Players struct {
	p map[Socket]model.Player
}

func (p *Players) Connect()