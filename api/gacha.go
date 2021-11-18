package api

import (
	"encoding/json"
	"fmt"
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

type GachaAPI struct {
	Auth *auth.Auth
	DB *database.GormDatabase
	Tx *gmtoken.GmtokenTx
}

type DrawGachaResponse struct {
	Results []model.GachaResult `json:"results"`
}

// localhost:8080/gacha/drawでガチャを引いて、キャラクターを取得
// -H "x-token:yyy"でトークン情報を受け取り、認証
// -d {"gacha_id":n, "times":x}でどのガチャを引くか、ガチャを何回引くかの情報を受け取る
func (a *GachaAPI) DrawGacha(ctx *gin.Context) {
	userId, err := a.Auth.GetUserId(ctx)
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	body, err := ioutil.ReadAll(ctx.Request.Body)
	defer ctx.Request.Body.Close()
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	var drawingGacha model.DrawingGacha
	if success := successOrAbort(ctx, http.StatusBadRequest, json.Unmarshal(body, &drawingGacha)); !success {
		return
	}
	contains, err := a.gachaIdContains(drawingGacha.GachaID)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	if !contains {
		err := fmt.Errorf("gacha_id error")
		if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
			return
		}
	}
	// 0以下回だけガチャを引くことは出来ない
	if drawingGacha.Times <= 0 {
		err := fmt.Errorf("times error")
		if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
			return
		}
	}
	enoughBal, err := a.checkBalance(userId, drawingGacha.Times)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	if !enoughBal {
		err := fmt.Errorf("balance of GameToken not enough")
		if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
			return
		}
	}
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	user, err := a.DB.GetUser(userId)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	// drawingGacha.Times分だけゲームトークンを焼却
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.Tx.BurnGmtoken(drawingGacha.Times, user.PrivateKey)); !success {
		return
	}
	characterInfos, err := a.DB.GetCharacterInfos(drawingGacha.GachaID)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	gachaCharacterIdsDrawed := drawGachaCharacterIds(characterInfos, drawingGacha.Times)
	var result model.GachaResult
	var results []model.GachaResult
	var userCharacters []model.UserCharacter
	count := 0
	for _, gacha_character_id := range gachaCharacterIdsDrawed {
		characterInfo := getCharacterInfo(characterInfos, gacha_character_id)
		result = model.GachaResult{CharacterID: gacha_character_id, Name: characterInfo.CharacterName}
		results = append(results, result)
		userCharacterId, err := util.CreateUUId()
		if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
			return
		}
		userCharacter := model.UserCharacter{UserCharacterID: userCharacterId, UserID: userId, GachaCharacterID: gacha_character_id}
		userCharacters = append(userCharacters, userCharacter)
		count += 1
		if count == 10000 {
			//	INSERT INTO `user_characters` (`user_character_id`,`user_id`,`gacha_character_id`)
			//	VALUES ('eaaada0c-3815-4da2-b791-3447a816a3e0','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6a8a4e-0ed8-11ec-93f3-a0c58933fdce')
			//	, ... ,
			//	('ff1583af-3f60-43de-839c-68094286e11a','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6d0b6d-0ed8-11ec-93f3-a0c58933fdce')
			if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateUserCharacters(userCharacters)); !success {
				return
			}
			userCharacters = userCharacters[:0]
			count = 0
		}
	}
	if len(userCharacters) != 0 {
		//	INSERT INTO `user_characters` (`user_character_id`,`user_id`,`gacha_character_id`)
		//	VALUES ('98b27372-8806-4d33-950a-68625ed6d687','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6c0f26-0ed8-11ec-93f3-a0c58933fdce')
		if success := successOrAbort(ctx, http.StatusInternalServerError, a.DB.CreateUserCharacters(userCharacters)); !success {
			return
		}
	}
	drawGachaResponse := &DrawGachaResponse{
		Results: results,
	}
	ctx.JSON(http.StatusOK, drawGachaResponse)
	//	{"results":[
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Sun"},
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Venus"},
	//		...
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Pluto"}
	//	]}
	//	が返る
}
