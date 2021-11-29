package database

import (
	"game-api-gin/model"
)

// user_charactersテーブルに引数userCharactersの情報を新規追加
func (d *GormDatabase) CreateUserCharacters(userCharacters []model.UserCharacter) error {
	//	INSERT INTO `user_characters` (`user_character_id`,`user_id`,`gacha_character_id`)
	//	VALUES ('eaaada0c-3815-4da2-b791-3447a816a3e0','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6a8a4e-0ed8-11ec-93f3-a0c58933fdce')
	//	, ... ,
	//	('ff1583af-3f60-43de-839c-68094286e11a','c2f0d74b-0321-4f87-930f-8d85350ee6d4','7b6d0b6d-0ed8-11ec-93f3-a0c58933fdce')
	return d.DB.Create(&userCharacters).Error
}

// dbからキャラクターのgacha_character_id、名前、weightの情報を取得
// ガチャidが引数gacha_idのキャラクターに限る
func (d *GormDatabase) GetCharacterInfos(gacha_id int) ([]model.CharacterInfo, error) {
	var characterInfos []model.CharacterInfo
	//	SELECT gacha_characters.gacha_character_id, characters.character_name, rarities.weight
	//	FROM `gacha_characters`
	//	join characters
	//	on gacha_characters.character_id = characters.id
	//	join rarities
	//	on gacha_characters.rarity_id = rarities.id
	//	WHERE gacha_id = 1
	err := d.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Where("gacha_id = ?", gacha_id).Scan(&characterInfos).Error
	return characterInfos, err
}

// user_charactersテーブルからユーザIDが引数user_idのデータを取得
func (d *GormDatabase) GetUserCharacters(user_id string) ([]model.UserCharacter, error) {
	var userCharacters []model.UserCharacter
	// SELECT * FROM `user_characters`  WHERE (user_id = '703a0b0a-1541-487e-be5b-906e9541b021')
	err := d.DB.Where("user_id = ?", user_id).Find(&userCharacters).Error
	return userCharacters, err
}

// dbから全てのキャラクターのgacha_character_id、名前、weightの情報を取得
func (d *GormDatabase) GetAllCharacterInfos() ([]model.CharacterInfo, error) {
	var allCharacterInfos []model.CharacterInfo
	//	SELECT gacha_characters.gacha_character_id, characters.character_name, rarities.weight
	//	FROM `gacha_characters`
	//	join characters
	//	on gacha_characters.character_id = characters.id
	//	join rarities
	//	on gacha_characters.rarity_id = rarities.id
	err := d.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Scan(&allCharacterInfos).Error
	return allCharacterInfos, err
}