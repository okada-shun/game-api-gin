package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
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

func TestGachaSuite(t *testing.T) {
	suite.Run(t, new(GachaSuite))
}

type GachaSuite struct {
	suite.Suite
	gapi       *GachaAPI
	auth       *auth.Auth
	db         *database.GormDatabase
	tx         *gmtoken.GmtokenTx
	ctx        *gin.Context
	recorder   *httptest.ResponseRecorder
	token      string
	privateKey string
	userId     string
}

func (s *GachaSuite) BeforeTest(string, string) {
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

	userId, err := util.CreateUUId()
	assert.Nil(s.T(), err)
	s.userId = userId
	privateKey, err := util.CreatePrivateKey()
	assert.Nil(s.T(), err)
	s.privateKey = privateKey
	s.db.CreateUser(model.User{UserID: userId, Name: "mike", PrivateKey: privateKey})
	token, err := s.auth.CreateToken(userId)
	assert.Nil(s.T(), err)
	s.token = token
	s.gapi.Tx.MintGmtoken(100, privateKey)

	s.T().Log("BeforeTest!!")
}

func (s *GachaSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	s.T().Log("AfterTest!!")
}

func (s *GachaSuite) Test_DrawGacha() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	body, _ := ioutil.ReadAll(s.recorder.Body)
	var drawGachaResponse DrawGachaResponse
	json.Unmarshal(body, &drawGachaResponse)
	assert.Equal(s.T(), 10, len(drawGachaResponse.Results))
	var gachaResults []model.GachaResult
	for i := 1; i < 10; i++ {
		gachaResults = append(gachaResults, drawGachaResponse.Results[i])
	}
	var characterIds []string
	s.db.DB.Table("gacha_characters").Select("gacha_character_id").Where("gacha_id = ?", 1).Scan(&characterIds)
	for _, v := range gachaResults {
		contains := test.Contains(characterIds, v.CharacterID)
		assert.Equal(s.T(), true, contains)
	}
	var characterNames []string
	s.db.DB.Table("gacha_characters").Select("characters.character_name").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Where("gacha_id = ?", 1).Scan(&characterNames)
	for _, v := range gachaResults {
		contains := test.Contains(characterNames, v.Name)
		assert.Equal(s.T(), true, contains)
	}
	_, bal, err := s.tx.GetAddressBalance(s.privateKey)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 90, bal)
}

func (s *GachaSuite) Test_DrawGacha_WithoutEnoughEth() {
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 500, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByInvalidToken() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token+"InvalidToken")

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByInvalidGachaID() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":10,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByZeroTimes() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":0}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByMinusTimes() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":-10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByExcessiveTimes() {
	s.gapi.Tx.TransferEth(200000000000000000, s.privateKey)

	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":1000}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
