package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"

	config "local.packages/config"
	database "local.packages/database"
	router "local.packages/router"
)

func main() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())

	config := config.NewConfig()
	db, err := database.NewDatabase(config)
	if err != nil {
		panic(err)
	}
	db_sql, err := db.DB.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close(db_sql)
	r, err := router.CreateRouter(db, config)
	if err != nil {
		panic(err)
	}
	r.Run("localhost:8080")
}