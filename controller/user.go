package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	
	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	model "local.packages/model"
	_ "github.com/go-sql-driver/mysql"
)

type TokenResponse struct {
	Token string `json:"token"`
}

// localhost:8080/user/createでユーザ情報を作成
// -d {"name":"aaa"}で名前データを受け取る
// UUIDでユーザIDを生成する
// ユーザIDからjwtでトークンを作成し、トークンを返す
func (a *UserGachaAPI) CreateUser(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	var user model.User
	if success := successOrAbort(ctx, http.StatusBadRequest, json.Unmarshal(body, &user)); !success {
		return
	}
	userId, err := createUUId()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	user.UserID = userId
	// 新規ユーザの秘密鍵を生成
	privateKey, err := crypto.GenerateKey()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	user.PrivateKey = privateKeyHex
	// ゲームトークンを100だけ鋳造し、新規ユーザに付与
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.MintGmtoken(100, user.PrivateKey)); !success {
		return
	}
	//	INSERT INTO `users` (`user_id`,`name`,`private_key`)
	//	VALUES ('95daec2b-287c-4358-ba6f-5c29e1c3cbdf','aaa','6e7eada90afb7e84bf5b4498c6adaa2d4014904644637d5fb355266944fbf93a')
	a.DB.CreateUser(user)
	// ユーザIDの文字列からjwtでトークン作成
	token, err := a.createToken(userId)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	// token = "生成されたトークンの文字列"
	tokenResponse := &TokenResponse{
		Token: token,
	}
	ctx.JSON(http.StatusOK, tokenResponse)
	// {"token":"生成されたトークンの文字列"}が返る
}

// ユーザIDからjwtでトークンを作成
// 有効期限は24時間に設定
// jwtのペイロードにはユーザIDと有効期限の時刻を設定
func (a *UserGachaAPI) createToken(userID string) (string, error) {
	// HS256は256ビットのハッシュ値を生成するアルゴリズム
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	// ペイロードにユーザIDと有効期限の時刻(24時間後)を設定
	token.Claims = jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	// 秘密鍵を取得
	signBytes, err := ioutil.ReadFile(a.Idrsa)
	if err != nil {
		return "", err
	}
	// 秘密鍵で署名
	tokenString, err := token.SignedString(signBytes)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// UUIDを生成
func createUUId() (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uu := u.String()
	return uu, nil
}

type UserResponse struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	GmtokenBalance int    `json:"gmtoken_balance"`
}

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// トークンからユーザIDを取り出し、dbからそのユーザIDのユーザの名前と秘密鍵データを取り出す
// 秘密鍵からユーザアドレスを生成
// コントラクトからそのユーザアドレスのゲームトークン残高を取り出し、返す
func (a *UserGachaAPI) GetUser(ctx *gin.Context) {
	userId, err := a.getUserId(ctx)
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	user := a.DB.GetUser(userId)
	address, balance, err := a.getAddressBalance(user.PrivateKey)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	userResponse := &UserResponse{
		Name:           user.Name,
		Address:        address.String(),
		GmtokenBalance: balance,
	}
	ctx.JSON(http.StatusOK, userResponse)
	// {"name":"aaa","address":"0x7a242084216fC7810aAe02c6deE5D9092C6B8fb9","gmtoken_balance":40}が返る
	// 有効期限が切れると{"code":400,"message":"Token is expired"}が返る
}

// jwtトークンを認証する
func (a *UserGachaAPI) verifyToken(tokenString string) (*jwt.Token, error) {
	// 秘密鍵を取得
	signBytes, err := ioutil.ReadFile(a.Idrsa)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signBytes, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// トークンからユーザID情報を取り出し、返す
func (a *UserGachaAPI) getUserId(ctx *gin.Context) (string, error) {
	tokenString := ctx.Request.Header.Get("x-token")
	token, err := a.verifyToken(tokenString)
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	// claims = map[exp:1.629639808e+09 userId:bdd4056a-f113-424c-9951-1eaaaf853e5c]
	userId := claims["userId"].(string)
	return userId, nil
}

// 引数の秘密鍵hexkeyからアドレスを生成
// コントラクトからそのアドレスのゲームトークン残高を取り出す
// アドレスと残高を返す
func (a *UserGachaAPI) getAddressBalance(hexkey string) (common.Address, int, error) {
	address, err := convertKeyToAddress(hexkey)
	if err != nil {
		return common.Address{}, 0, err
	}
	bal, err := a.Gmtoken.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return common.Address{}, 0, err
	}
	balance, _ := strconv.Atoi(bal.String())
	return address, balance, nil
}

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// -d {"name":"bbb"}で更新する名前データを受け取る
// トークンからユーザIDを取り出し、dbからそのユーザIDのユーザの情報を更新
func (a *UserGachaAPI) UpdateUser(ctx *gin.Context) {
	userId, err := a.getUserId(ctx)
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	var user model.User
	if success := successOrAbort(ctx, http.StatusBadRequest, json.Unmarshal(body, &user)); !success {
		return
	}
	// dbでnameとaddressを更新
	// UPDATE `users` SET `name`='bbb' WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	a.DB.UpdateUser(user, userId)
	ctx.JSON(http.StatusOK, nil)
}

