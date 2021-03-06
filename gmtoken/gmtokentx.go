package gmtoken

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"

	"game-api-gin/config"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

type GmtokenTx struct {
	Config    *config.Config
	Ethclient *ethclient.Client
	Gmtoken   *Gmtoken
}

// Gmtokenのインスタンスを返す
func getGmtokenInstance(config *config.Config) (*Gmtoken, error) {
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
	instance, err := NewGmtoken(contractAddress, client)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// イーサリアムネットワークに接続するクライアントを返す
func getEthclient(config *config.Config) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Ethereum.NetworkURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// GmtokenTxインスタンス作成
func NewGmtokenTx(config *config.Config) (*GmtokenTx, error) {
	gmtokenInstance, err := getGmtokenInstance(config)
	if err != nil {
		return nil, err
	}
	ethclient, err := getEthclient(config)
	if err != nil {
		return nil, err
	}
	return &GmtokenTx{
		Config:    config,
		Gmtoken:   gmtokenInstance,
		Ethclient: ethclient,
	}, nil
}

// 16進数の秘密鍵文字列をイーサリアムアドレスに変換
func ConvertKeyToAddress(hexkey string) (common.Address, error) {
	// 16進数の秘密鍵文字列を読み込む
	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return common.Address{}, err
	}
	// 秘密鍵から公開鍵を生成
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	// 公開鍵からアドレスを生成
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

// 引数の秘密鍵hexkeyからアドレスを生成
// コントラクトからそのアドレスのゲームトークン残高を取り出す
// アドレスと残高を返す
func (g *GmtokenTx) GetAddressBalance(hexkey string) (common.Address, int, error) {
	address, err := ConvertKeyToAddress(hexkey)
	if err != nil {
		return common.Address{}, 0, err
	}
	bal, err := g.Gmtoken.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return common.Address{}, 0, err
	}
	balance, _ := strconv.Atoi(bal.String())
	return address, balance, nil
}

// 引数valだけethを転送する
// 引数hexkeyの秘密鍵から生成されるアドレスに転送する
// トランザクションの送信者は、minter_private_key.txtの秘密鍵から生成されるアドレスである
func (g *GmtokenTx) TransferEth(val int64, hexkey string) error {
	// 16進数の秘密鍵文字列をアドレスに変換
	address, err := ConvertKeyToAddress(hexkey)
	if err != nil {
		return err
	}
	// トランザクションを送るアドレスの秘密鍵を読み込む
	privateKeyBytes, err := ioutil.ReadFile(g.Config.Ethereum.MinterPrivateKey)
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(string(privateKeyBytes))
	if err != nil {
		return err
	}
	// 秘密鍵からアドレスを生成
	fromAddress, err := ConvertKeyToAddress(string(privateKeyBytes))
	if err != nil {
		return err
	}
	// ナンスを生成
	nonce, err := g.Ethclient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	// 転送するイーサの量を設定(ここではval eth)
	value := big.NewInt(val)
	// ガス価格を設定（SuggestGasPriceで平均のガス価格を取得）
	gasPrice, err := g.Ethclient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(gasPrice) // 20000000000
	var gasLimit uint64 = 10000000
	var data []byte
	// ナンス、転送先アドレス、転送するイーサ量、ガス制限、ガス価格、データからトランザクションを作成
	tx := types.NewTransaction(nonce, address, value, gasLimit, gasPrice, data)
	// チェーンID(ネットワークID)を取得
	chainID, err := g.Ethclient.NetworkID(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(chainID) // 5777
	// 送信者の秘密鍵を使用してトランザクションに署名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	// fmt.Println(signedTx)
	// トランザクションを送信
	err = g.Ethclient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	// fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
	return nil
}

// コントラクトから、引数valだけゲームトークンを鋳造する
// 鋳造されたゲームトークンは、引数hexkeyの秘密鍵から生成されるアドレスに付与される
// トランザクションの送信者は、minter_private_key.txtの秘密鍵から生成されるアドレスである
func (g *GmtokenTx) MintGmtoken(val int, hexkey string) error {
	// 16進数の秘密鍵文字列をアドレスに変換
	address, err := ConvertKeyToAddress(hexkey)
	if err != nil {
		return err
	}
	// トランザクションを送るアドレスの秘密鍵を読み込む
	privateKeyBytes, err := ioutil.ReadFile(g.Config.Ethereum.MinterPrivateKey)
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(string(privateKeyBytes))
	if err != nil {
		return err
	}
	// 秘密鍵からアドレスを生成
	fromAddress, err := ConvertKeyToAddress(string(privateKeyBytes))
	if err != nil {
		return err
	}
	// ナンスを生成
	nonce, err := g.Ethclient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	// 転送するイーサの量を設定(ここでは0)
	value := big.NewInt(0)
	// ガス価格を設定（SuggestGasPriceで平均のガス価格を取得）
	gasPrice, err := g.Ethclient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(gasPrice) // 20000000000
	// GameTokenコントラクトのアドレスを読み込む
	contractAddressBytes, err := ioutil.ReadFile(g.Config.Ethereum.ContractAddress)
	if err != nil {
		return err
	}
	contractAddress := common.HexToAddress(string(contractAddressBytes))
	// スマートコントラクトのmint関数
	mintFnSignature := []byte("mint(address,uint256)")
	// Keccak256ハッシュを生成
	hash := sha3.NewLegacyKeccak256()
	hash.Write(mintFnSignature)
	methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0x40c10f19
	paddedAddress := common.LeftPadBytes(address.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAddress)) // 0x00000000000000000000000001cf71979d4a11b0f15ca36b315cf9e77561c070
	// 転送するゲームトークンの量を設定
	amount := new(big.Int)
	amount.SetString(strconv.Itoa(val), 10)
	// 32バイト幅まで左側を0で埋める
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAmount)) // 0x0000000000000000000000000000000000000000000000000000000000000064
	// メソッドIDと32バイト幅トークン量を連結
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	/*
		// EstimateGasでガス制限を推定
		gasLimit, err := g.Ethclient.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		})
		if err != nil {
			return err
		}
	*/
	var gasLimit uint64 = 10000000
	// ナンス、コントラクトのアドレス、転送するイーサ量(0)、ガス制限、ガス価格、データからトランザクションを作成
	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	// チェーンID(ネットワークID)を取得
	chainID, err := g.Ethclient.NetworkID(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(chainID) // 5777
	// 送信者の秘密鍵を使用してトランザクションに署名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	// fmt.Println(signedTx) // &{0xc00006e5a0 {13858186735851446500 49580235901 0xef4fa0} {<nil>} {<nil>} {<nil>}}
	// トランザクションを送信
	err = g.Ethclient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	// fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x8369c729025e98fd73e01c6e99724bb397bc58274b963b6eab75f1bd10dc39a1
	return nil
}

// コントラクトから、引数valだけゲームトークンを焼却する
// 引数hexkeyの秘密鍵から生成されるアドレスの持つゲームトークンを焼却する
// トランザクションの送信者は、引数hexkeyの秘密鍵から生成されるアドレスである
func (g *GmtokenTx) BurnGmtoken(val int, hexkey string) error {
	// 16進数の秘密鍵文字列を読み込む
	privateKey, err := crypto.HexToECDSA(hexkey)
	if err != nil {
		return err
	}
	// 16進数の秘密鍵文字列をアドレスに変換
	address, err := ConvertKeyToAddress(hexkey)
	if err != nil {
		return err
	}
	// トランザクションを送るアドレス
	fromAddress := address
	// ナンスを生成
	nonce, err := g.Ethclient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}
	// 転送するイーサの量を設定(ここでは0)
	value := big.NewInt(0)
	// ガス価格を設定（SuggestGasPriceで平均のガス価格を取得）
	gasPrice, err := g.Ethclient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(gasPrice) // 20000000000
	// GameTokenコントラクトのアドレスを読み込む
	contractAddressBytes, err := ioutil.ReadFile(g.Config.Ethereum.ContractAddress)
	if err != nil {
		return err
	}
	contractAddress := common.HexToAddress(string(contractAddressBytes))
	// スマートコントラクトのburn関数
	burnFnSignature := []byte("burn(uint256)")
	// Keccak256ハッシュを生成
	hash := sha3.NewLegacyKeccak256()
	hash.Write(burnFnSignature)
	methodID := hash.Sum(nil)[:4]
	// fmt.Println(hexutil.Encode(methodID)) // 0x42966c68
	// 転送するゲームトークンの量を設定
	amount := new(big.Int)
	amount.SetString(strconv.Itoa(val), 10)
	// 32バイト幅まで左側を0で埋める
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	// fmt.Println(hexutil.Encode(paddedAmount)) // 0x000000000000000000000000000000000000000000000000000000000000000a
	// メソッドIDと32バイト幅トークン量を連結
	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	/*
		// EstimateGasでガス制限を推定
		gasLimit, err := g.Ethclient.EstimateGas(context.Background(), ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		})
		if err != nil {
			return err
		}
	*/
	var gasLimit uint64 = 10000000
	// ナンス、コントラクトのアドレス、転送するイーサ量(0)、ガス制限、ガス価格、データからトランザクションを作成
	tx := types.NewTransaction(nonce, contractAddress, value, gasLimit, gasPrice, data)
	// チェーンID(ネットワークID)を取得
	chainID, err := g.Ethclient.NetworkID(context.Background())
	if err != nil {
		return err
	}
	// fmt.Println(chainID) // 5777
	// 送信者の秘密鍵を使用してトランザクションに署名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	// fmt.Println(signedTx) // &{0xc00012bf20 {13858187368593859236 638888717301 0xef4fa0} {<nil>} {<nil>} {<nil>}}
	// トランザクションを送信
	err = g.Ethclient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	// fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0xf98c12a353eceacafe606397493d0d321628f1a70bb147697d1539a2a9ca9199
	return nil
}
