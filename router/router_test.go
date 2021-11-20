package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"game-api-gin/api"
	"game-api-gin/auth"
	"game-api-gin/config"
	"game-api-gin/database"
	"game-api-gin/gmtoken"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	client = &http.Client{}
)

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

type RouterSuite struct {
	suite.Suite
	a *auth.Auth
	db *database.GormDatabase
	tx *gmtoken.GmtokenTx
	server *httptest.Server
}

func (s *RouterSuite) BeforeTest(string, string) {
	config, err := config.NewConfig()
	assert.Nil(s.T(), err)
	a := auth.NewAuth(config)
	s.a = a
	d, err := database.NewDatabase(config)
	assert.Nil(s.T(), err)
	s.db = d
	t, err := gmtoken.NewGmtokenTx(config)
	assert.Nil(s.T(), err)
	s.tx = t
	r := CreateRouter(s.a, s.db, s.tx)
	s.server = httptest.NewServer(r)
}

func (s *RouterSuite) AfterTest(string, string) {
	s.db.Close()
	s.server.Close()
}

func (s *RouterSuite) TestRouter() {
	req := s.newRequest("GET", "", "")
	res, err := client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	req = s.newRequest("POST", "user/create", `{"name":"testName"}`)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	createUserResponse := api.CreateUserResponse{}
	json.NewDecoder(res.Body).Decode(&createUserResponse)
	token := createUserResponse.Token

	req = s.newRequest("GET", "user/get", "")
	req.Header.Add("x-token", token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	getUserResponse := api.GetUserResponse{}
	json.NewDecoder(res.Body).Decode(&getUserResponse)
	assert.Equal(s.T(), "testName", getUserResponse.Name)
	assert.Equal(s.T(), 100, getUserResponse.GmtokenBalance)

	req = s.newRequest("PUT", "user/update", `{"name":"updateName"}`)
	req.Header.Add("x-token", token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	req = s.newRequest("GET", "user/get", "")
	req.Header.Add("x-token", token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	getUserResponse = api.GetUserResponse{}
	json.NewDecoder(res.Body).Decode(&getUserResponse)
	assert.Equal(s.T(), "updateName", getUserResponse.Name)

	//s.tx.TransferEth(200000000000000000, getUserResponse.Address)

	req = s.newRequest("POST", "gacha/draw", `{"gacha_id":1,"times":10}`)
	req.Header.Add("x-token", token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)

	req = s.newRequest("GET", "character/list", "")
	req.Header.Add("x-token", token)
	res, err = client.Do(req)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, res.StatusCode)
}

func (s *RouterSuite) newRequest(method, url, jsonStr string) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", s.server.URL, url), bytes.NewBuffer([]byte(jsonStr)))
	assert.Nil(s.T(), err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("accept", "application/json")
	return req
}

/*
func (s *RouterSuite) TestCreateRouterFail() {
	config := &config.Config{}
	a := &auth.Auth{}
	d := &database.GormDatabase{}
	t := &gmtoken.GmtokenTx{}
	r := CreateRouter(a, d, t)
	assert.Error(s.T(), err)
}
*/
