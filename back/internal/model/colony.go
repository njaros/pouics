package model

import (
	"fmt"
	"log"
	mathtools "pouic/internal/helpers/math_tools"
)

type Colony struct {
    Id string       `json:"id"`
    PlayerId string `json:"player_id"`
    MaxSize int     `json:"max_size"`
    PosX float64    `json:"pos_x"`
    PosY float64    `json:"pos_y"`
    TargetX float64 `json:"target_x"`
    TargetY float64 `json:"target_y"`
    Speed float64   `json:"speed"`
}

func (c *Colony) Move() {
    if c.Speed != 0 {
        logString := fmt.Sprintf("colonie %s (addr %p) moved from (%f, %f) ", c.Id, &c, c.PosX, c.PosY)
        dirX, dirY := mathtools.PointsToVector(c.PosX, c.PosY, c.TargetX, c.TargetY)
        addX, addY := mathtools.NormDecomposition(c.Speed, dirX, dirY)
        if addX >= dirX {
            c.PosX = c.TargetX
            c.PosY = c.TargetY
            c.Speed = 0
        } else {
            c.PosX += addX
            c.PosY += addY
        }
        logString += fmt.Sprintf("to (%f, %f)\n", c.PosX, c.PosY)
        log.Print(logString)
    }
}