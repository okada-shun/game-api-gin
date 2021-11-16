package api

import (
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	config "local.packages/config"
	gmtoken "local.packages/gmtoken"
)

// Gmtokenのインスタンスを返す
func NewGmtoken(config *config.Config) (*gmtoken.Gmtoken, error) {
	client, err := ethclient.Dial(config.NetworkURL)
	if err != nil {
		return nil, err
	}
	// GameTokenコントラクトのアドレスを読み込む
	contractAddressBytes, err := ioutil.ReadFile(config.ContractAddress)
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
	client, err := ethclient.Dial(config.NetworkURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}
