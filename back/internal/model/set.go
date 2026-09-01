package model

type set struct {
    Id int `json:"id"`
    Player_id int `json:"player_id"`
    MaxSize int `json:"max_size"`
    PosX int `json:"pos_x"`
    PosY int `json:"pos_y"`
    TargetX int `json:"target_x"`
    TargetY int `json:"target_y"`
    Speed float32 `json:"speed"`
}