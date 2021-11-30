package auth

import (
	"io/ioutil"

	"game-api-gin/config"

	"github.com/dgrijalva/jwt-go"
)

type Auth struct {
	Idrsa string
}

func NewAuth(config *config.Config) *Auth {
	return &Auth{
		Idrsa: config.Auth.Idrsa,
	}
}

// jwtトークンを照合する
func (a *Auth) VerifyToken(tokenString string) (*jwt.Token, error) {
	// 秘密鍵を取得
	signBytes, err := ioutil.ReadFile(a.Idrsa)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signBytes, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
