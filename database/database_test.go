package database

import (
	"testing"

	"game-api-gin/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

type DatabaseSuite struct {
	suite.Suite
	db *GormDatabase
}

func (s *DatabaseSuite) BeforeTest(suiteName, testName string) {
	config, err := config.NewConfig()
	assert.Nil(s.T(), err)
	d, err := NewDatabase(config)
	assert.Nil(s.T(), err)
	s.db = d
}

func (s *DatabaseSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func TestNewDatabase(t *testing.T) {
	config, err := config.NewConfig()
	assert.Nil(t, err)
	d, err := NewDatabase(config)
	assert.Nil(t, err)
	d.Close()
}

func TestNewDatabase_PasswordFail(t *testing.T) {
	config, err := config.NewConfig()
	assert.Nil(t, err)
	config.Mysql.Password = "asdf"
	_, err = NewDatabase(config)
	assert.Error(t, err)
}

func TestNewDatabase_UserFail(t *testing.T) {
	config, err := config.NewConfig()
	assert.Nil(t, err)
	config.Mysql.User = "asdf"
	_, err = NewDatabase(config)
	assert.Error(t, err)
}
