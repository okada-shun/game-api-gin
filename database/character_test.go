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
	var expectedCharacterInfos01 [10]model.CharacterInfo
	expectedCharacterInfos01[0] =
		model.CharacterInfo{GachaCharacterID: "0d390f6b-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 1}
	expectedCharacterInfos01[1] =
		model.CharacterInfo{GachaCharacterID: "0d39f724-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 5}
	expectedCharacterInfos01[2] =
		model.CharacterInfo{GachaCharacterID: "0d3aeef4-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 5}
	expectedCharacterInfos01[3] =
		model.CharacterInfo{GachaCharacterID: "0d3c07d8-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 5}
	expectedCharacterInfos01[4] =
		model.CharacterInfo{GachaCharacterID: "0d3cf0ad-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 14}
	expectedCharacterInfos01[5] =
		model.CharacterInfo{GachaCharacterID: "0d3df5e3-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 14}
	expectedCharacterInfos01[6] =
		model.CharacterInfo{GachaCharacterID: "0d3ecd91-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 14}
	expectedCharacterInfos01[7] =
		model.CharacterInfo{GachaCharacterID: "0d3fea83-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 14}
	expectedCharacterInfos01[8] =
		model.CharacterInfo{GachaCharacterID: "0d40d0ae-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 14}
	expectedCharacterInfos01[9] =
		model.CharacterInfo{GachaCharacterID: "0d41b476-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 14}
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
	var expectedCharacterInfos02 [10]model.CharacterInfo
	expectedCharacterInfos02[0] =
		model.CharacterInfo{GachaCharacterID: "0d42a9a3-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 14}
	expectedCharacterInfos02[1] =
		model.CharacterInfo{GachaCharacterID: "0d437909-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 14}
	expectedCharacterInfos02[2] =
		model.CharacterInfo{GachaCharacterID: "0d445545-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 14}
	expectedCharacterInfos02[3] =
		model.CharacterInfo{GachaCharacterID: "0d453701-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 1}
	expectedCharacterInfos02[4] =
		model.CharacterInfo{GachaCharacterID: "0d460f79-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 5}
	expectedCharacterInfos02[5] =
		model.CharacterInfo{GachaCharacterID: "0d474188-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 5}
	expectedCharacterInfos02[6] =
		model.CharacterInfo{GachaCharacterID: "0d4826b4-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 5}
	expectedCharacterInfos02[7] =
		model.CharacterInfo{GachaCharacterID: "0d48f4d8-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 14}
	expectedCharacterInfos02[8] =
		model.CharacterInfo{GachaCharacterID: "0d49ed9c-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 14}
	expectedCharacterInfos02[9] =
		model.CharacterInfo{GachaCharacterID: "0d4abc80-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 14}
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
	var expectedCharacterInfos03 [10]model.CharacterInfo
	expectedCharacterInfos03[0] =
		model.CharacterInfo{GachaCharacterID: "0d4b8cc9-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 14}
	expectedCharacterInfos03[1] =
		model.CharacterInfo{GachaCharacterID: "0d4cc775-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 14}
	expectedCharacterInfos03[2] =
		model.CharacterInfo{GachaCharacterID: "0d4d8d92-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 14}
	expectedCharacterInfos03[3] =
		model.CharacterInfo{GachaCharacterID: "0d4f0347-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 14}
	expectedCharacterInfos03[4] =
		model.CharacterInfo{GachaCharacterID: "0d4fc8b3-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 14}
	expectedCharacterInfos03[5] =
		model.CharacterInfo{GachaCharacterID: "0d509e02-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 14}
	expectedCharacterInfos03[6] =
		model.CharacterInfo{GachaCharacterID: "0d51acef-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 1}
	expectedCharacterInfos03[7] =
		model.CharacterInfo{GachaCharacterID: "0d527e11-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 5}
	expectedCharacterInfos03[8] =
		model.CharacterInfo{GachaCharacterID: "0d5362ad-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 5}
	expectedCharacterInfos03[9] =
		model.CharacterInfo{GachaCharacterID: "0d544ad9-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 5}
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
	var expectedAllCharacterInfos [30]model.CharacterInfo
	expectedAllCharacterInfos[0] =
		model.CharacterInfo{GachaCharacterID: "0d390f6b-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 1}
	expectedAllCharacterInfos[1] =
		model.CharacterInfo{GachaCharacterID: "0d39f724-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 5}
	expectedAllCharacterInfos[2] =
		model.CharacterInfo{GachaCharacterID: "0d3aeef4-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 5}
	expectedAllCharacterInfos[3] =
		model.CharacterInfo{GachaCharacterID: "0d3c07d8-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 5}
	expectedAllCharacterInfos[4] =
		model.CharacterInfo{GachaCharacterID: "0d3cf0ad-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 14}
	expectedAllCharacterInfos[5] =
		model.CharacterInfo{GachaCharacterID: "0d3df5e3-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 14}
	expectedAllCharacterInfos[6] =
		model.CharacterInfo{GachaCharacterID: "0d3ecd91-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 14}
	expectedAllCharacterInfos[7] =
		model.CharacterInfo{GachaCharacterID: "0d3fea83-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 14}
	expectedAllCharacterInfos[8] =
		model.CharacterInfo{GachaCharacterID: "0d40d0ae-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 14}
	expectedAllCharacterInfos[9] =
		model.CharacterInfo{GachaCharacterID: "0d41b476-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 14}
	expectedAllCharacterInfos[10] =
		model.CharacterInfo{GachaCharacterID: "0d42a9a3-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 14}
	expectedAllCharacterInfos[11] =
		model.CharacterInfo{GachaCharacterID: "0d437909-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 14}
	expectedAllCharacterInfos[12] =
		model.CharacterInfo{GachaCharacterID: "0d445545-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 14}
	expectedAllCharacterInfos[13] =
		model.CharacterInfo{GachaCharacterID: "0d453701-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 1}
	expectedAllCharacterInfos[14] =
		model.CharacterInfo{GachaCharacterID: "0d460f79-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 5}
	expectedAllCharacterInfos[15] =
		model.CharacterInfo{GachaCharacterID: "0d474188-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 5}
	expectedAllCharacterInfos[16] =
		model.CharacterInfo{GachaCharacterID: "0d4826b4-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 5}
	expectedAllCharacterInfos[17] =
		model.CharacterInfo{GachaCharacterID: "0d48f4d8-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 14}
	expectedAllCharacterInfos[18] =
		model.CharacterInfo{GachaCharacterID: "0d49ed9c-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 14}
	expectedAllCharacterInfos[19] =
		model.CharacterInfo{GachaCharacterID: "0d4abc80-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 14}
	expectedAllCharacterInfos[20] =
		model.CharacterInfo{GachaCharacterID: "0d4b8cc9-5197-11ec-830e-a0c58933fdce", CharacterName: "Mercury", Weight: 14}
	expectedAllCharacterInfos[21] =
		model.CharacterInfo{GachaCharacterID: "0d4cc775-5197-11ec-830e-a0c58933fdce", CharacterName: "Venus", Weight: 14}
	expectedAllCharacterInfos[22] =
		model.CharacterInfo{GachaCharacterID: "0d4d8d92-5197-11ec-830e-a0c58933fdce", CharacterName: "Earth", Weight: 14}
	expectedAllCharacterInfos[23] =
		model.CharacterInfo{GachaCharacterID: "0d4f0347-5197-11ec-830e-a0c58933fdce", CharacterName: "Mars", Weight: 14}
	expectedAllCharacterInfos[24] =
		model.CharacterInfo{GachaCharacterID: "0d4fc8b3-5197-11ec-830e-a0c58933fdce", CharacterName: "Jupiter", Weight: 14}
	expectedAllCharacterInfos[25] =
		model.CharacterInfo{GachaCharacterID: "0d509e02-5197-11ec-830e-a0c58933fdce", CharacterName: "Saturn", Weight: 14}
	expectedAllCharacterInfos[26] =
		model.CharacterInfo{GachaCharacterID: "0d51acef-5197-11ec-830e-a0c58933fdce", CharacterName: "Uranus", Weight: 1}
	expectedAllCharacterInfos[27] =
		model.CharacterInfo{GachaCharacterID: "0d527e11-5197-11ec-830e-a0c58933fdce", CharacterName: "Neptune", Weight: 5}
	expectedAllCharacterInfos[28] =
		model.CharacterInfo{GachaCharacterID: "0d5362ad-5197-11ec-830e-a0c58933fdce", CharacterName: "Pluto", Weight: 5}
	expectedAllCharacterInfos[29] =
		model.CharacterInfo{GachaCharacterID: "0d544ad9-5197-11ec-830e-a0c58933fdce", CharacterName: "Sun", Weight: 5}
	expectedSet := mapset.NewSet()
	for _, v := range expectedAllCharacterInfos {
		expectedSet.Add(v)
	}
	assert.Equal(s.T(), expectedSet, actualSet)
}
