package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	"game-api-gin/config"
	"game-api-gin/database"
	"game-api-gin/router"
)

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	db, err := database.NewDatabase(config)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r, err := router.CreateRouter(db, config)
	if err != nil {
		panic(err)
	}
	r.Run("localhost:8080")
}