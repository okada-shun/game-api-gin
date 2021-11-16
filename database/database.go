package database

import (
	"database/sql"
	"io/ioutil"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	config "local.packages/config"
)

type GormDatabase struct {
	DB *gorm.DB
}

// databaseインスタンスを返す
func NewDatabase(config *config.Config) (*GormDatabase, error) {
	passwordBytes, err := ioutil.ReadFile(config.MysqlPass)
	if err != nil {
		return nil, err
	}
	userBytes, err := ioutil.ReadFile(config.MysqlUser)
	if err != nil {
		return nil, err
	}
	dsn := string(userBytes)+":"+string(passwordBytes)+"@/game_user?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Logger = db.Logger.LogMode(logger.Info)
	return &GormDatabase{
		DB: db,
	}, nil
}

func (d *GormDatabase) Close(db_sql *sql.DB) {
	db_sql.Close()
}
