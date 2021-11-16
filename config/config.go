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
		Idrsa: "/home/okada_shun/.ssh/id_rsa",
		MinterPrivateKey: "/home/okada_shun/game-api-gin/gmtoken/minter_private_key.txt",
		ContractAddress: "/home/okada_shun/game-api-gin/gmtoken/GameToken_address.txt",
		NetworkURL: "ws://localhost:7545",
		MysqlPass: "/home/okada_shun/.ssh/mysql_password",
		MysqlUser: "/home/okada_shun/.ssh/mysql_user",
	}
}
