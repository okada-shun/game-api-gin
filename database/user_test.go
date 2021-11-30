package database

import (
	"game-api-gin/model"
	"game-api-gin/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DatabaseSuite) TestUser() {
	user, err := s.db.GetUser("asdf")
	assert.Error(s.T(), err)

	userId, err := util.CreateUUId()
	require.NoError(s.T(), err)
	privateKey, err := util.CreatePrivateKey()
	require.NoError(s.T(), err)

	alice := model.User{UserID: userId, Name: "alice", PrivateKey: privateKey}
	s.db.CreateUser(alice)
	user, err = s.db.GetUser(userId)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), alice, user)

	alice.Name = "aliceUPDATED"
	s.db.UpdateUser(alice, userId)
	user, err = s.db.GetUser(userId)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), model.User{UserID: userId, Name: "aliceUPDATED", PrivateKey: privateKey}, user)
}
