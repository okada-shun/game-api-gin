package database

import (
	"game-api-gin/model"
	"game-api-gin/util"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (s *DatabaseSuite) TestUserCharacters() {
	userId, err := util.CreateUUId()
	require.NoError(s.T(), err)
	var userCharactersCreated []model.UserCharacter
	for i := 0; i < 10; i++ {
		userCharacterID, _ := util.CreateUUId()
		gachaCharacterID, _ := util.CreateUUId()
		userCharacter := model.UserCharacter{UserCharacterID: userCharacterID, UserID: userId, GachaCharacterID: gachaCharacterID}
		userCharactersCreated = append(userCharactersCreated, userCharacter)
	}
	s.db.CreateUserCharacters(userCharactersCreated)
	setCreated := mapset.NewSet()
	for _, v := range userCharactersCreated {
		setCreated.Add(v)
	}
	userCharactersGot, err := s.db.GetUserCharacters(userId)
	require.NoError(s.T(), err)
	setGot := mapset.NewSet()
	for _, v := range userCharactersGot {
		setGot.Add(v)
	}
	assert.Equal(s.T(), setCreated, setGot)
}

func (s *DatabaseSuite) TestGetCharacterInfos() {
	actualCharacterInfos01, err := s.db.GetCharacterInfos(1)
	require.NoError(s.T(), err)
	actualSet01 := mapset.NewSet()
	for _, v := range actualCharacterInfos01 {
		actualSet01.Add(v)
	}
	var expectedCharacterInfos01 []model.CharacterInfo
	s.db.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Where("gacha_id = ?", 1).Scan(&expectedCharacterInfos01)
	expectedSet01 := mapset.NewSet()
	for _, v := range expectedCharacterInfos01 {
		expectedSet01.Add(v)
	}
	assert.Equal(s.T(), expectedSet01, actualSet01)

	actualCharacterInfos02, err := s.db.GetCharacterInfos(2)
	require.NoError(s.T(), err)
	actualSet02 := mapset.NewSet()
	for _, v := range actualCharacterInfos02 {
		actualSet02.Add(v)
	}
	var expectedCharacterInfos02 []model.CharacterInfo
	s.db.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Where("gacha_id = ?", 2).Scan(&expectedCharacterInfos02)
	expectedSet02 := mapset.NewSet()
	for _, v := range expectedCharacterInfos02 {
		expectedSet02.Add(v)
	}
	assert.Equal(s.T(), expectedSet02, actualSet02)

	actualCharacterInfos03, err := s.db.GetCharacterInfos(3)
	require.NoError(s.T(), err)
	actualSet03 := mapset.NewSet()
	for _, v := range actualCharacterInfos03 {
		actualSet03.Add(v)
	}
	var expectedCharacterInfos03 []model.CharacterInfo
	s.db.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").
		Where("gacha_id = ?", 3).Scan(&expectedCharacterInfos03)
	expectedSet03 := mapset.NewSet()
	for _, v := range expectedCharacterInfos03 {
		expectedSet03.Add(v)
	}
	assert.Equal(s.T(), expectedSet03, actualSet03)
}

func (s *DatabaseSuite) TestGetAllCharacters() {
	actualAllCharacterInfos, err := s.db.GetAllCharacterInfos()
	require.NoError(s.T(), err)
	actualSet := mapset.NewSet()
	for _, v := range actualAllCharacterInfos {
		actualSet.Add(v)
	}
	var expectedAllCharacterInfos []model.CharacterInfo
	s.db.DB.Table("gacha_characters").Select("gacha_characters.gacha_character_id, characters.character_name, rarities.weight").
		Joins("join characters on gacha_characters.character_id = characters.id").
		Joins("join rarities on gacha_characters.rarity_id = rarities.id").Scan(&expectedAllCharacterInfos)
	expectedSet := mapset.NewSet()
	for _, v := range expectedAllCharacterInfos {
		expectedSet.Add(v)
	}
	assert.Equal(s.T(), expectedSet, actualSet)
}
