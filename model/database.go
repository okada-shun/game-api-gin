package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	MysqlPass string
	MysqlUser string
}

// databaseインスタンスを返す
func NewDatabase(mysqlpass string, mysqluser string) *Database {
	return &Database{
		MysqlPass: mysqlpass,
		MysqlUser: mysqluser,
	}
}

// DB(game_user)からコネクション取得
func (d *Database) getConnection() (*gorm.DB, error) {
	passwordBytes, err := ioutil.ReadFile(d.MysqlPass)
	if err != nil {
		return nil, err
	}
	userBytes, err := ioutil.ReadFile(d.MysqlUser)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(string(userBytes)+":"+string(passwordBytes)+"@/game_user?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Logger = db.Logger.LogMode(logger.Info)
	return db, nil
}

type User struct {
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	PrivateKey string `json:"private_key"`
}

type Character struct {
	GachaCharacterID string `json:"gacha_character_id"`
	CharacterName    string `json:"character_name"`
	Weight           uint   `json:"weight"`
}

type UserCharacter struct {
	UserCharacterID  string `json:"user_character_id"`
	UserID           string `json:"user_id"`
	GachaCharacterID string `json:"gacha_character_id"`
}

// usersテーブルにユーザ情報を新規追加
func (d *Database) CreateUser(user User) {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	//	INSERT INTO `users` (`user_id`,`name`,`private_key`)
	//	VALUES ('95daec2b-287c-4358-ba6f-5c29e1c3cbdf','aaa','6e7eada90afb7e84bf5b4498c6adaa2d4014904644637d5fb355266944fbf93a')
	db.Create(&user)
}

// usersテーブルからユーザIDが引数userIdのユーザの情報を取得
func (d *Database) GetUser(userId string) User {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	var user User
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	db.Where("user_id = ?", userId).Find(&user)
	return user
}

// usersテーブルからユーザIDが引数userIdのユーザの情報を、引数userのものに更新
func (d *Database) UpdateUser(user User, userId string) {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	// UPDATE `users` SET `name`='bbb' WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	db.Model(&user).Where("user_id = ?", userId).Update("name", user.Name)
}

// user_charactersテーブルに引数userCharactersの情報を新規追加
func (d *Database) CreateUserCharacters(userCharacters []UserCharacter) {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	//	INSERT INTO `user_characters` (`user_character_id`,`user_id`,`gacha_character_id`)
	//	VALUES ('eaaada0c-3815-4da2-b791-3447a816a3e0','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6a8a4e-0ed8-11ec-93f3-a0c58933fdce')
	//	, ... ,
	//	('ff1583af-3f60-43de-839c-68094286e11a','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6d0b6d-0ed8-11ec-93f3-a0c58933fdce')
	db.Create(&userCharacters)
}

// gacha_charactersテーブルからガチャIDを全て取得
func (d *Database) GetGachaIds() []int {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	var gachaIds []int
	// SELECT gacha_id FROM `gacha_characters`
	db.Table("gacha_characters").Select("gacha_id").Scan(&gachaIds)
	return gachaIds
}

// dbからキャラクターのgacha_character_id、名前、weightの情報を取得
// ガチャidが引数gacha_idのキャラクターに限る
func (d *Database) GetCharacters(gacha_id int) []Character {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	var charactersList []Character
	//	SELECT gacha_characters.gacha_character_id, characters.character_name, rarities.weight
	//	FROM `gacha_characters`
	//	join characters
	//	on gacha_characters.character_id = characters.id
	//	join rarities
	//	on gacha_characters.rarity_id = rarities.id
	//	WHERE gacha_id = 1
	db.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Where("gacha_id = ?", gacha_id).Scan(&charactersList)
	return charactersList
}

// user_charactersテーブルからユーザIDが引数user_idのデータを取得
func (d *Database) GetUserCharacters(user_id string) []UserCharacter {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	var userCharactersList []UserCharacter
	// SELECT * FROM `user_characters`  WHERE (user_id = '703a0b0a-1541-487e-be5b-906e9541b021')
	db.Where("user_id = ?", user_id).Find(&userCharactersList)
	return userCharactersList
}

// dbから全てのキャラクターのgacha_character_id、名前、weightの情報を取得
func (d *Database) GetAllCharacters() []Character {
	db, err := d.getConnection()
	if err != nil {
		fmt.Println(err.Error(), http.StatusInternalServerError)
	}
	db_sql, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer db_sql.Close()
	var charactersList []Character
	//	SELECT gacha_characters.gacha_character_id, characters.character_name, rarities.weight
	//	FROM `gacha_characters`
	//	join characters
	//	on gacha_characters.character_id = characters.id
	//	join rarities
	//	on gacha_characters.rarity_id = rarities.id
	db.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Scan(&charactersList)
	return charactersList
}