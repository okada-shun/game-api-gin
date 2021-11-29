package model

type User struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	PrivateKey string `json:"private_key"`
}
