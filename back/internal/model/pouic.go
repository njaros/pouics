package model

type Pouic struct {
    Id int                  `json:"id"`
    SetId int              `json:"set_id"`
    Name string             `json:"name"`
    Health int              `json:"health"`
    Happyness int           `json:"happyness"`
    Stress int              `json:"stress"`
    Stats [3]int            `json:"stats"`
    Brain [64][4]float32    `json:"brain"`    
}