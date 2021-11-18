package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"game-api-gin/auth"
	"game-api-gin/database"
	"game-api-gin/gmtoken"
	"game-api-gin/model"
	"game-api-gin/util"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type UserAPI struct {
	Auth *auth.Auth
	DB *database.GormDatabase
	Tx *gmtoken.GmtokenTx
}

type CreateUserRes struct {
	Token string `json:"token"`
}

type GetUserRes struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	GmtokenBalance int    `json:"gmtoken_balance"`
}

// localhost:8080/user/createでユーザ情報を作成
// -d {"name":"aaa"}で名前データを受け取る
// UUIDでユーザIDを生成する
// ユーザIDからjwtでトークンを作成し、トークンを返す
func (a *UserAPI) CreateUser(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	var user model.User
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, json.Unmarshal(body, &user)); !success {
		return
	}
	userId, err := util.CreateUUId()
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	user.UserID = userId
	// 新規ユーザの秘密鍵を生成
	privateKeyHex, err := util.CreatePrivateKey()
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	user.PrivateKey = privateKeyHex
	// ゲームトークンを100だけ鋳造し、新規ユーザに付与
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, a.Tx.MintGmtoken(100, user.PrivateKey)); !success {
		return
	}
	//	INSERT INTO `users` (`user_id`,`name`,`private_key`)
	//	VALUES ('95daec2b-287c-4358-ba6f-5c29e1c3cbdf','aaa','6e7eada90afb7e84bf5b4498c6adaa2d4014904644637d5fb355266944fbf93a')
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateUser(user)); !success {
		return
	}
	// ユーザIDの文字列からjwtでトークン作成
	token, err := a.Auth.CreateToken(userId)
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	// token = "生成されたトークンの文字列"
	createUserRes := &CreateUserRes{
		Token: token,
	}
	ctx.JSON(http.StatusOK, createUserRes)
	// {"token":"生成されたトークンの文字列"}が返る
}

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// トークンからユーザIDを取り出し、dbからそのユーザIDのユーザの名前と秘密鍵データを取り出す
// 秘密鍵からユーザアドレスを生成
// コントラクトからそのユーザアドレスのゲームトークン残高を取り出し、返す
func (a *UserAPI) GetUser(ctx *gin.Context) {
	userId, err := a.Auth.GetUserId(ctx)
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	user, err := a.DB.GetUser(userId)
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	address, balance, err := a.Tx.GetAddressBalance(user.PrivateKey)
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	getUserRes := &GetUserRes{
		Name:           user.Name,
		Address:        address.String(),
		GmtokenBalance: balance,
	}
	ctx.JSON(http.StatusOK, getUserRes)
	// {"name":"aaa","address":"0x7a242084216fC7810aAe02c6deE5D9092C6B8fb9","gmtoken_balance":40}が返る
	// 有効期限が切れると{"code":400,"message":"Token is expired"}が返る
}

// -H "x-token:yyy"でトークン情報を受け取り、ユーザ認証
// -d {"name":"bbb"}で更新する名前データを受け取る
// トークンからユーザIDを取り出し、dbからそのユーザIDのユーザの情報を更新
func (a *UserAPI) UpdateUser(ctx *gin.Context) {
	userId, err := a.Auth.GetUserId(ctx)
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	var user model.User
	if success := util.SuccessOrAbort(ctx, http.StatusBadRequest, json.Unmarshal(body, &user)); !success {
		return
	}
	// dbでnameとaddressを更新
	// UPDATE `users` SET `name`='bbb' WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	if success := util.SuccessOrAbort(ctx, http.StatusInternalServerError, a.DB.UpdateUser(user, userId)); !success {
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
