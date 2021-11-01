package config

// Idrsa: jwtトークンの作成・認証に使用するサーバーの秘密鍵へのパス
// MinterPrivateKey: MintGmtoken関数で使用する、Minterの秘密鍵へのパス
// ContractAddress: GameTokenコントラクトのアドレスへのパス
// NetworkURL: イーサリアムのプロバイダーURL
// MysqlPass: MySQLのパスワードへのパス
// MysqlUser: MySQLのユーザ名へのパス
type Config struct {
	Idrsa string
	MinterPrivateKey string
	ContractAddress string
	NetworkURL string
	MysqlPass string
	MysqlUser string
}

// configインスタンスを作成
func NewConfig() *Config {
	return &Config{
		Idrsa: "../.ssh/id_rsa",
		MinterPrivateKey: "./gmtoken/minter_private_key.txt",
		ContractAddress: "./gmtoken/GameToken_address.txt",
		NetworkURL: "ws://localhost:7545",
		MysqlPass: "../.ssh/mysql_password",
		MysqlUser: "../.ssh/mysql_user",
	}
}
