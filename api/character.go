package api

import (
	"net/http"

	"game-api-gin/auth"
	"game-api-gin/database"
	"game-api-gin/model"
	
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CharacterAPI struct {
	Auth *auth.Auth
	DB *database.GormDatabase
}

type GetCharacterListResponse struct {
	Characters []model.Character `json:"characters"`
}

// localhost:8080/character/listでユーザが所持しているキャラクター一覧情報を取得
// -H "x-token:yyy"でトークン情報を受け取り、認証
func (a *CharacterAPI) GetCharacterList(ctx *gin.Context) {
	userId, err := a.Auth.GetUserId(ctx)
	if success := successOrAbort(ctx, http.StatusBadRequest, err); !success {
		return
	}
	allCharacterInfos, err := a.DB.GetAllCharacterInfos()
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	userCharacters, err := a.DB.GetUserCharacters(userId)
	if success := successOrAbort(ctx, http.StatusInternalServerError, err); !success {
		return
	}
	var characters []model.Character
	var character model.Character
	if len(userCharacters) == 0 {
		characters = make([]model.Character, 0)
	} else {
		for _, v := range userCharacters {
			gacha_character_id := v.GachaCharacterID
			characterInfo := getCharacterInfo(allCharacterInfos, gacha_character_id)
			characterName := characterInfo.CharacterName
			character = model.Character{UserCharacterID: v.UserCharacterID, CharacterID: gacha_character_id, Name: characterName}
			characters = append(characters, character)
		}
	}
	getCharacterListResponse := &GetCharacterListResponse{
		Characters: characters,
	}
	ctx.JSON(http.StatusOK, getCharacterListResponse)
	//	{"characters":[
	//		{"userCharacterID":"02091c4d-1011-4615-8fbb-fd9e681153d4","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Sun"},
	//		{"userCharacterID":"0fed4c04-153c-4980-9a66-1424f1f7a445","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Venus"},
	//		...
	//		{"userCharacterID":"95a281d5-86f0-4251-a4cb-5873231f4a96","characterID":"c115174c-05ad-11ec-8679-a0c58933fdce","name":"Pluto"}
	//	]}
	//	が返る
}
