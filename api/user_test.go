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
	"game-api-gin/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

type UserSuite struct {
	suite.Suite
	uapi       *UserAPI
	auth       *auth.Auth
	db         *database.GormDatabase
	tx         *gmtoken.GmtokenTx
	ctx        *gin.Context
	recorder   *httptest.ResponseRecorder
	token      string
	privateKey string
	userId     string
}

func (s *UserSuite) BeforeTest(string, string) {
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

	s.T().Log("BeforeTest!!")
}

func (s *UserSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
	s.T().Log("AfterTest!!")
}

func (s *UserSuite) Test_HelloWorld() {
	s.ctx.Request = httptest.NewRequest("GET", "/", nil)
	assert.Equal(s.T(), 200, s.recorder.Code)
}

func (s *UserSuite) Test_CreateUser() {
	s.ctx.Request = httptest.NewRequest("POST", "/user/create", strings.NewReader(`{"name":"tom"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.uapi.CreateUser(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	body, _ := ioutil.ReadAll(s.recorder.Body)
	var createUserResponse CreateUserResponse
	json.Unmarshal(body, &createUserResponse)
	token, err := s.auth.VerifyToken(createUserResponse.Token)
	assert.Nil(s.T(), err)
	claims := token.Claims.(jwt.MapClaims)
	userId := claims["userId"].(string)
	var privateKey string
	s.db.DB.Table("users").Select("private_key").Where("user_id = ?", userId).Scan(&privateKey)
	_, bal, err := s.tx.GetAddressBalance(privateKey)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 100, bal)
}

func (s *UserSuite) Test_CreateUser_ByIntName() {
	s.ctx.Request = httptest.NewRequest("POST", "/user/create", strings.NewReader(`{"name":123}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.uapi.CreateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_CreateUser_ByNil() {
	s.ctx.Request = httptest.NewRequest("POST", "/user/create", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")

	s.uapi.CreateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_GetUser() {
	s.ctx.Request = httptest.NewRequest("GET", "/user/get", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.uapi.GetUser(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)

	body, _ := ioutil.ReadAll(s.recorder.Body)
	var getUserResponse GetUserResponse
	json.Unmarshal(body, &getUserResponse)
	address, err := gmtoken.ConvertKeyToAddress(s.privateKey)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "mike", getUserResponse.Name)
	assert.Equal(s.T(), address.Hex(), getUserResponse.Address)
}

func (s *UserSuite) Test_GetUser_ByInvalidToken() {
	s.ctx.Request = httptest.NewRequest("GET", "/user/get", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token+"InvalidToken")

	s.uapi.GetUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUser() {
	s.ctx.Request = httptest.NewRequest("PUT", "/user/update", strings.NewReader(`{"name":"wang"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.uapi.UpdateUser(s.ctx)

	assert.Equal(s.T(), 200, s.recorder.Code)
	user, err := s.db.GetUser(s.userId)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "wang", user.Name)
}

func (s *UserSuite) Test_UpdateUser_ByInvalidToken() {
	s.ctx.Request = httptest.NewRequest("PUT", "/user/update", strings.NewReader(`{"name":"wang"}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token+"InvalidToken")

	s.uapi.UpdateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUser_ByIntName() {
	s.ctx.Request = httptest.NewRequest("PUT", "/user/update", strings.NewReader(`{"name":123}`))
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.uapi.UpdateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}

func (s *UserSuite) Test_UpdateUser_ByNil() {
	s.ctx.Request = httptest.NewRequest("PUT", "/user/update", nil)
	s.ctx.Request.Header.Set("Content-Type", "application/json")
	s.ctx.Request.Header.Set("x-token", s.token)

	s.uapi.UpdateUser(s.ctx)

	assert.Equal(s.T(), 400, s.recorder.Code)
}
