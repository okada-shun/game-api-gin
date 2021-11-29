package router

import (
	"net/http"

	"game-api-gin/api"
	"game-api-gin/auth"
	"game-api-gin/database"
	"game-api-gin/gmtoken"

	"github.com/gin-gonic/gin"
)

func CreateRouter(auth *auth.Auth, db *database.GormDatabase, gmtokenTx *gmtoken.GmtokenTx) *gin.Engine {
	router := gin.Default()
	userHandler := &api.UserAPI{
		Auth: auth,
		DB:   db,
		Tx:   gmtokenTx,
	}
	gachaHandler := &api.GachaAPI{
		Auth: auth,
		DB:   db,
		Tx:   gmtokenTx,
	}
	characterHandler := &api.CharacterAPI{
		Auth: auth,
		DB:   db,
	}
	router.GET("/", home)
	router.POST("/user/create", userHandler.CreateUser)
	router.GET("/user/get", userHandler.GetUser)
	router.PUT("/user/update", userHandler.UpdateUser)
	router.POST("/gacha/draw", gachaHandler.DrawGacha)
	router.GET("/character/list", characterHandler.GetCharacterList)
	return router
}

// {"message":"Hello World"}をlocalhost:8080画面に表示
func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
