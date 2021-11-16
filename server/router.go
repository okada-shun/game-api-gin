package server
 
import (
	"net/http"
	
  "github.com/gin-gonic/gin"
	controller "local.packages/controller"
)
 
func GetRouter() *gin.Engine {
	router := gin.Default()
	// userGachaAPIインスタンスを作成
	userGachaAPI := controller.NewUserGachaAPI()
	router = startServer(router, userGachaAPI)
	return router
}

// {"message":"Hello World"}をlocalhost:8080画面に表示
func home(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

// サーバー起動
func startServer(router *gin.Engine, userGachaAPI *controller.UserGachaAPI) *gin.Engine {
	router.GET("/", home)
	router.POST("/user/create", userGachaAPI.CreateUser)
	router.GET("/user/get", userGachaAPI.GetUser)
	router.PUT("/user/update", userGachaAPI.UpdateUser)
	router.POST("/gacha/draw", userGachaAPI.DrawGacha)
	router.GET("/character/list", userGachaAPI.GetCharacterList)
	return router
}