package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	server "local.packages/server"
)

func main() {
	// 乱数のシード値を設定
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	// サーバー起動
	r := server.GetRouter()
	r.Run("localhost:8080")
}