package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	wr "github.com/mroth/weightedrand"
	"game-api-gin/database"
	"game-api-gin/gmtoken"
	"game-api-gin/model"
	"game-api-gin/util"
)

type GachaAPI struct {
	Idrsa string
	MinterPrivateKey string
	ContractAddress string
	Gmtoken *gmtoken.Gmtoken
	DB *database.GormDatabase
	Ethclient *ethclient.Client
	AuthToken *AuthTokenAPI
	Transaction *TransactionAPI
}

// localhost:8080/gacha/drawでガチャを引いて、キャラクターを取得
// -H "x-token:yyy"でトークン情報を受け取り、認証
// -d {"gacha_id":n, "times":x}でどのガチャを引くか、ガチャを何回引くかの情報を受け取る
func (a *GachaAPI) DrawGacha(ctx *gin.Context) {
	userId, err := a.AuthToken.GetUserId(ctx)
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
	if success := successOrAbort(ctx, http.StatusInternalServerError, a.Transaction.burnGmtoken(drawingGacha.Times, user.PrivateKey)); !success {
		return
	}
	charactersList, err := a.DB.GetCharacters(drawingGacha.GachaID)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	gachaCharacterIdsDrawed := drawGachaCharacterIds(charactersList, drawingGacha.Times)
	var characterInfo model.CharacterResponse
	var results []model.CharacterResponse
	var userCharacters []model.UserCharacter
	count := 0
	for _, gacha_character_id := range gachaCharacterIdsDrawed {
		character := getCharacterInfo(charactersList, gacha_character_id)
		characterInfo = model.CharacterResponse{CharacterID: gacha_character_id, Name: character.CharacterName}
		results = append(results, characterInfo)
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
	resultResponse := &model.ResultResponse{
		Results: results,
	}
	ctx.JSON(http.StatusOK, resultResponse)
	//	{"results":[
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Sun"},
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Venus"},
	//		...
	//		{"characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Pluto"}
	//	]}
	//	が返る
}

// dbのgacha_charactersテーブルからgacha_id一覧を取得
// 引数のgachaIdがその中に含まれていたらtrue、含まれていなかったらfalseを返す
func (a *GachaAPI) gachaIdContains(gachaId int) (bool, error) {
	// SELECT gacha_id FROM `gacha_characters`
	gachaIds, err := a.DB.GetGachaIds()
	if err != nil {
		return false, err
	}
	for _, v := range gachaIds {
		if v == gachaId {
			return true, nil
		}
	}
	return false, nil
}

// dbのusersテーブルからuser_idが引数userIdのユーザ情報を取得
// コントラクトからそのユーザアドレスのゲームトークン残高を取得
// 引数のtimesが残高以下だったらtrue、残高より大きかったらfalseを返す
func (a *GachaAPI) checkBalance(userId string, times int) (bool, error) {
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	user, err := a.DB.GetUser(userId)
	if err != nil {
		return false, err
	}
	_, balance, err := a.Transaction.getAddressBalance(user.PrivateKey)
	if err != nil {
		return false, err
	}
	return times <= balance, nil
}

// charactersListからキャラクターのgacha_character_idとweightを取り出しchoicesに格納
// times回分だけchoicesからWeighted Random Selectionを実行
func drawGachaCharacterIds(charactersList []model.Character, times int) []string {
	var choices []wr.Choice
	for i := 0; i < len(charactersList); i++ {
		choices = append(choices, wr.Choice{Item: charactersList[i].GachaCharacterID, Weight: charactersList[i].Weight})
	}
	var gachaCharacterIdsDrawed []string
	for i := 0; i < times; i++ {
		chooser, _ := wr.NewChooser(choices...)
		gachaCharacterIdsDrawed = append(gachaCharacterIdsDrawed, chooser.Pick().(string))
	}
	return gachaCharacterIdsDrawed
}

// 引数のcharactersListからGachaCharacterIDが引数gacha_character_idのデータを取得
func getCharacterInfo(charactersList []model.Character, gacha_character_id string) model.Character {
	for i := 0; i < len(charactersList); i++ {
		if charactersList[i].GachaCharacterID == gacha_character_id {
			return charactersList[i]
		}
	}
	return model.Character{}
}

// localhost:8080/character/listでユーザが所持しているキャラクター一覧情報を取得
// -H "x-token:yyy"でトークン情報を受け取り、認証
func (a *GachaAPI) GetCharacterList(ctx *gin.Context) {
	userId, err := a.AuthToken.GetUserId(ctx)
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	allCharactersList, err := a.DB.GetAllCharacters()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	userCharactersList, err := a.DB.GetUserCharacters(userId)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	var characters []model.UserCharacterResponse
	var userCharacterInfo model.UserCharacterResponse
	if len(userCharactersList) == 0 {
		characters = make([]model.UserCharacterResponse, 0)
	} else {
		for _, v := range userCharactersList {
			gacha_character_id := v.GachaCharacterID
			character := getCharacterInfo(allCharactersList, gacha_character_id)
			characterName := character.CharacterName
			userCharacterInfo = model.UserCharacterResponse{UserCharacterID: v.UserCharacterID, CharacterID: gacha_character_id, Name: characterName}
			characters = append(characters, userCharacterInfo)
		}
	}
	charactersResponse := &model.CharactersResponse{
		Characters: characters,
	}
	ctx.JSON(http.StatusOK, charactersResponse)
	//	{"characters":[
	//		{"userCharacterID":"02091c4d-1011-4615-8fbb-fd9e681153d4","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Sun"},
	//		{"userCharacterID":"0fed4c04-153c-4980-9a66-1424f1f7a445","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Venus"},
	//		...
	//		{"userCharacterID":"95a281d5-86f0-4251-a4cb-5873231f4a96","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Pluto"}
	//	]}
	//	が返る
}

