package model

type DrawingGacha struct {
	GachaID int `json:"gacha_id"`
	Times   int `json:"times"`
}

type GachaResult struct {
	CharacterID string `json:"characterID"`
	Name        string `json:"name"`
}
