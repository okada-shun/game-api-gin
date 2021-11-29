package api

import (
	"game-api-gin/model"
	wr "github.com/mroth/weightedrand"
)

// dbのgacha_charactersテーブルからgacha_id一覧を取得
// 引数のgachaIdがその中に含まれていたらtrue、含まれていなかったらfalseを返す
func (a *GachaAPI) gachaIdContains(gachaId int) (bool, error) {
	// SELECT gacha_id FROM `gacha_characters`
	gachaIds, err := a.DB.GetGachaIds()
	if err != nil {
		return false, err
	}
	for _, v := range gachaIds {
		if v == gachaId {
			return true, nil
		}
	}
	return false, nil
}

// dbのusersテーブルからuser_idが引数userIdのユーザ情報を取得
// コントラクトからそのユーザアドレスのゲームトークン残高を取得
// 引数のtimesが残高以下だったらtrue、残高より大きかったらfalseを返す
func (a *GachaAPI) checkBalance(userId string, times int) (bool, error) {
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	user, err := a.DB.GetUser(userId)
	if err != nil {
		return false, err
	}
	_, balance, err := a.Tx.GetAddressBalance(user.PrivateKey)
	if err != nil {
		return false, err
	}
	return times <= balance, nil
}

// characterInfosからキャラクターのgacha_character_idとweightを取り出しchoicesに格納
// times回分だけchoicesからWeighted Random Selectionを実行
func drawGachaCharacterIds(characterInfos []model.CharacterInfo, times int) []string {
	var choices []wr.Choice
	for i := 0; i < len(characterInfos); i++ {
		choices = append(choices, wr.Choice{Item: characterInfos[i].GachaCharacterID, Weight: characterInfos[i].Weight})
	}
	var gachaCharacterIdsDrawed []string
	for i := 0; i < times; i++ {
		chooser, _ := wr.NewChooser(choices...)
		gachaCharacterIdsDrawed = append(gachaCharacterIdsDrawed, chooser.Pick().(string))
	}
	return gachaCharacterIdsDrawed
}

// 引数のcharacterInfosからGachaCharacterIDが引数gacha_character_idのデータを取得
func getCharacterInfo(characterInfos []model.CharacterInfo, gacha_character_id string) model.CharacterInfo {
	for i := 0; i < len(characterInfos); i++ {
		if characterInfos[i].GachaCharacterID == gacha_character_id {
			return characterInfos[i]
		}
	}
	return model.CharacterInfo{}
}
