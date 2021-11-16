package model

type User struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	PrivateKey string `json:"private_key"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	GmtokenBalance int    `json:"gmtoken_balance"`
}
