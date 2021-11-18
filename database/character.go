package database

import (
	"game-api-gin/model"
)

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