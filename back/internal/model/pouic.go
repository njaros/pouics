package model

type Pouic struct {
    Id string               `json:"id"`
    ColonyId string         `json:"set_id"`
    PlayerId string         `json:"player_id"`
    Name string             `json:"name"`
    Health int              `json:"health"`
    Happyness int           `json:"happyness"`
    Stress int              `json:"stress"`
    Stats [3]int            `json:"stats"`
    Brain [64][4]float32    `json:"brain"`    
}