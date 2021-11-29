package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// トークンからユーザID情報を取り出し、返す
func (a *Auth) GetUserId(ctx *gin.Context) (string, error) {
	tokenString := ctx.Request.Header.Get("x-token")
	token, err := a.verifyToken(tokenString)
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	// claims = map[exp:1.629639808e+09 userId:bdd4056a-f113-424c-9951-1eaaaf853e5c]
	userId := claims["userId"].(string)
	return userId, nil
}
