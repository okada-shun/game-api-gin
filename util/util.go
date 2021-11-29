package util

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// UUIDを生成
func CreateUUId() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uu := u.String()
	return uu, nil
}

// 秘密鍵を生成
func CreatePrivateKey() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	return privateKeyHex, nil
}
