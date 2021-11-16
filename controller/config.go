package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	config "local.packages/config"
	model "local.packages/model"
	gmtoken "local.packages/gmtoken"
)

// Idrsa: jwtトークンの作成・認証に使用するサーバーの秘密鍵へのパス
// MinterPrivateKey: MintGmtoken関数で使用する、Minterの秘密鍵へのパス
// ContractAddress: GameTokenコントラクトのアドレスへのパス
type UserGachaAPI struct {
	Idrsa string
	MinterPrivateKey string
	ContractAddress string
	Gmtoken *gmtoken.Gmtoken
	DB *model.Database
	Ethclient *ethclient.Client
}

// userGachaAPIインスタンスを作成
func NewUserGachaAPI() *UserGachaAPI {
	config := config.NewConfig()
	return &UserGachaAPI{
		Idrsa: config.Idrsa,
		MinterPrivateKey: config.MinterPrivateKey,
		ContractAddress: config.ContractAddress,
		Gmtoken: newGmtoken(config.NetworkURL, config.ContractAddress),
		DB: model.NewDatabase(config.MysqlPass, config.MysqlUser),
		Ethclient: newEthclient(config.NetworkURL),
	}
}

// Gmtokenのインスタンスを返す
func newGmtoken(url string, address string) *gmtoken.Gmtoken {
	gmtokenInstance, err := getGmtokenInstance(url, address)
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	return gmtokenInstance
}

// イーサリアムネットワークに接続するクライアントを返す
func newEthclient(url string) *ethclient.Client {
	client, err := getEthclient(url)
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	return client
}

// Gmtokenのインスタンスを取得
func getGmtokenInstance(url string, address string) (*gmtoken.Gmtoken, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	// GameTokenコントラクトのアドレスを読み込む
	contractAddressBytes, err := ioutil.ReadFile(address)
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

// イーサリアムネットワークに接続するクライアントを取得
func getEthclient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}