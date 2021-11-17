package api

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"game-api-gin/config"
	"game-api-gin/gmtoken"
)

// Gmtokenのインスタンスを返す
func NewGmtoken(config *config.Config) (*gmtoken.Gmtoken, error) {
	client, err := ethclient.Dial(config.Ethereum.NetworkURL)
	if err != nil {
		return nil, err
	}
	// GameTokenコントラクトのアドレスを読み込む
	contractAddressBytes, err := ioutil.ReadFile(config.Ethereum.ContractAddress)
	if err != nil {
		return nil, err
	}
	contractAddress := common.HexToAddress(string(contractAddressBytes))
	// 上のコントラクトアドレスのGmtokenインスタンスを作成
	instance, err := gmtoken.NewGmtoken(contractAddress, client)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// イーサリアムネットワークに接続するクライアントを返す
func NewEthclient(config *config.Config) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Ethereum.NetworkURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}
