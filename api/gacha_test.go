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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestGachaSuite(t *testing.T) {
	suite.Run(t, new(GachaSuite))
}

type GachaSuite struct {
	suite.Suite
	uapi *UserAPI
	gapi *GachaAPI
	auth *auth.Auth
	db *database.GormDatabase
	tx *gmtoken.GmtokenTx
	ctx *gin.Context
	recorder *httptest.ResponseRecorder
	token string
	privatekey string
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
	s.uapi = &UserAPI{Auth: s.auth, DB: s.db, Tx: s.tx}
	s.gapi = &GachaAPI{Auth: s.auth, DB: s.db, Tx: s.tx}
	s.T().Log("BeforeTest!!")
}

func (s *GachaSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	s.T().Log("AfterTest!!")
}

func (s *GachaSuite) Test_DrawGacha() {
	s.ctx.Request = httptest.NewRequest("POST", "/user/create", strings.NewReader(`{"name":"satou"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.uapi.CreateUser(s.ctx)
	body, _ := ioutil.ReadAll(s.recorder.Body)
	var createUserResponse CreateUserResponse
	json.Unmarshal(body, &createUserResponse)
	s.token = createUserResponse.Token
	s.ctx.Request.Header.Set("x-token", s.token)

	userId, err := s.auth.GetUserId(s.ctx)
	assert.Nil(s.T(), err)
	user, err := s.db.GetUser(userId)
	assert.Nil(s.T(), err)
	s.privatekey = user.PrivateKey

	s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	body, _ = ioutil.ReadAll(s.recorder.Body)
	var drawGachaResponse DrawGachaResponse
	json.Unmarshal(body, &drawGachaResponse)
	assert.Equal(s.T(), 10, len(drawGachaResponse.Results))
}

func (s *GachaSuite) Test_DrawGacha_ByInvalidToken() {
	//s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", "InvalidToken")

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByInvalidGachaID() {
	//s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":10,"times":10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByZeroTimes() {
	//s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":0}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByMinusTimes() {
	//s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":-10}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *GachaSuite) Test_DrawGacha_ByExcessiveTimes() {
	//s.gapi.Tx.TransferEth(200000000000000000, s.privatekey)
	s.ctx.Request = httptest.NewRequest("POST", "/gacha/draw", strings.NewReader(`{"gacha_id":1,"times":1000}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.gapi.DrawGacha(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
