package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"game-api-gin/auth"
	"game-api-gin/config"
	"game-api-gin/database"
	"game-api-gin/gmtoken"
	"game-api-gin/router"
)

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	auth := auth.NewAuth(config)
	db, err := database.NewDatabase(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	gmtokenTx, err := gmtoken.NewGmtokenTx(config)
	if err != nil {
		panic(err)
	}
	r := router.CreateRouter(auth, db, gmtokenTx)
	r.Run("localhost:8080")
}