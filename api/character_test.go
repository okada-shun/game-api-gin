package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"game-api-gin/auth"
	"game-api-gin/config"
	"game-api-gin/database"
	"game-api-gin/gmtoken"
	"game-api-gin/model"
	"game-api-gin/test"
	"game-api-gin/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCharacterSuite(t *testing.T) {
	suite.Run(t, new(CharacterSuite))
}

type CharacterSuite struct {
	suite.Suite
	gapi     *GachaAPI
	capi     *CharacterAPI
	auth     *auth.Auth
	db       *database.GormDatabase
	tx       *gmtoken.GmtokenTx
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
	token    string
	userId   string
}

func (s *CharacterSuite) BeforeTest(string, string) {
	s.recorder = httptest.NewRecorder()
	s.ctx, _ = gin.CreateTestContext(s.recorder)
	config, err := config.NewConfig()
	assert.Nil(s.T(), err)
	a := auth.NewAuth(config)
	s.auth = a
	d, err := database.NewDatabase(config)
	assert.Nil(s.T(), err)
	s.db = d
	t, err := gmtoken.NewGmtokenTx(config)
	assert.Nil(s.T(), err)
	s.tx = t
	s.gapi = &GachaAPI{Auth: s.auth, DB: s.db, Tx: s.tx}
	s.capi = &CharacterAPI{Auth: s.auth, DB: s.db}

	userId, err := util.CreateUUId()
	assert.Nil(s.T(), err)
	s.userId = userId
	privateKey, err := util.CreatePrivateKey()
	assert.Nil(s.T(), err)
	s.db.CreateUser(model.User{UserID: userId, Name: "mike", PrivateKey: privateKey})
	token, err := s.auth.CreateToken(userId)
	assert.Nil(s.T(), err)
	s.token = token

	s.T().Log("BeforeTest!!")
}

func (s *CharacterSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	s.T().Log("AfterTest!!")
}

func (s *CharacterSuite) Test_GetCharacterList() {
	s.ctx.Request = httptest.NewRequest("GET", "/character/list", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.capi.GetCharacterList(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	body, _ := ioutil.ReadAll(s.recorder.Body)
	var getCharacterListResponse GetCharacterListResponse
	json.Unmarshal(body, &getCharacterListResponse)
	assert.Equal(s.T(), 0, len(getCharacterListResponse.Characters))
	assert.Equal(s.T(), []model.Character{}, getCharacterListResponse.Characters)
}

func (s *CharacterSuite) Test_GetCharacterList_WithCharacter() {
	var gachaCharacterIds []string
	s.db.DB.Table("gacha_characters").Select("gacha_character_id").Where("gacha_id = ?", 1).Scan(&gachaCharacterIds)
	var userCharacters []model.UserCharacter
	var userCharacterIds []string
	for _, v := range gachaCharacterIds {
		userCharacterId, err := util.CreateUUId()
		assert.Nil(s.T(), err)
		userCharacterIds = append(userCharacterIds, userCharacterId)
		userCharacter := model.UserCharacter{UserCharacterID: userCharacterId, UserID: s.userId, GachaCharacterID: v}
		userCharacters = append(userCharacters, userCharacter)
	}
	s.gapi.DB.CreateUserCharacters(userCharacters)

	s.ctx.Request = httptest.NewRequest("GET", "/character/list", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.capi.GetCharacterList(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	body, _ := ioutil.ReadAll(s.recorder.Body)
	var getCharacterListResponse GetCharacterListResponse
	json.Unmarshal(body, &getCharacterListResponse)
	assert.Equal(s.T(), 10, len(getCharacterListResponse.Characters))
	for _, v := range getCharacterListResponse.Characters {
		contains := test.Contains(userCharacterIds, v.UserCharacterID)
		assert.Equal(s.T(), true, contains)
	}
	var characterIds []string
	s.db.DB.Table("gacha_characters").Select("gacha_character_id").Where("gacha_id = ?", 1).Scan(&characterIds)
	for _, v := range getCharacterListResponse.Characters {
		contains := test.Contains(characterIds, v.CharacterID)
		assert.Equal(s.T(), true, contains)
	}
	var characterNames []string
	s.db.DB.Table("gacha_characters").Select("characters.character_name").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Where("gacha_id = ?", 1).Scan(&characterNames)
	for _, v := range getCharacterListResponse.Characters {
		contains := test.Contains(characterNames, v.Name)
		assert.Equal(s.T(), true, contains)
	}
}

func (s *CharacterSuite) Test_GetCharacterList_ByInvalidToken() {
	var gachaCharacterIds []string
	s.db.DB.Table("gacha_characters").Select("gacha_character_id").Where("gacha_id = ?", 1).Scan(&gachaCharacterIds)
	var userCharacters []model.UserCharacter
	for _, v := range gachaCharacterIds {
		userCharacterId, err := util.CreateUUId()
		assert.Nil(s.T(), err)
		userCharacter := model.UserCharacter{UserCharacterID: userCharacterId, UserID: s.userId, GachaCharacterID: v}
		userCharacters = append(userCharacters, userCharacter)
	}
	s.gapi.DB.CreateUserCharacters(userCharacters)

	s.ctx.Request = httptest.NewRequest("GET", "/character/list", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token+"InvalidToken")

	s.capi.GetCharacterList(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
