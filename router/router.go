package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"game-api-gin/api"
	"game-api-gin/config"
	"game-api-gin/database"
)

func CreateRouter(db *database.GormDatabase, config *config.Config) (*gin.Engine, error) {
	router := gin.Default()
	gmtoken, err := api.NewGmtoken(config)
	if err != nil {
		return nil, err
	}
	ethclient, err := api.NewEthclient(config)
	if err != nil {
		return nil, err
	}
	authToken := api.NewAuthToken(config)
	transaction, err := api.NewTransaction(config)
	if err != nil {
		return nil, err
	}
	userHandler := &api.UserAPI{
		Idrsa: config.AuthToken.Idrsa,
		MinterPrivateKey: config.Ethereum.MinterPrivateKey,
		ContractAddress: config.Ethereum.ContractAddress,
		Gmtoken: gmtoken,
		DB: db,
		Ethclient: ethclient,
		AuthToken: authToken,
		Transaction: transaction,
	}
	gachaHandler := &api.GachaAPI{
		Idrsa: config.AuthToken.Idrsa,
		MinterPrivateKey: config.Ethereum.MinterPrivateKey,
		ContractAddress: config.Ethereum.ContractAddress,
		Gmtoken: gmtoken,
		DB: db,
		Ethclient: ethclient,
		AuthToken: authToken,
		Transaction: transaction,
	}
	router.GET("/", home)
	router.POST("/user/create", userHandler.CreateUser)
	router.GET("/user/get", userHandler.GetUser)
	router.PUT("/user/update", userHandler.UpdateUser)
	router.POST("/gacha/draw", gachaHandler.DrawGacha)
	router.GET("/character/list", gachaHandler.GetCharacterList)
	return router, nil
}

// {"message":"Hello World"}をlocalhost:8080画面に表示
func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}
