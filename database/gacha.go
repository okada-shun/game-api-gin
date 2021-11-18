package database

// gacha_charactersテーブルからガチャIDを全て取得
func (d *GormDatabase) GetGachaIds() ([]int, error) {
	var gachaIds []int
	// SELECT gacha_id FROM `gacha_characters`
	err := d.DB.Table("gacha_characters").Select("gacha_id").Scan(&gachaIds).Error
	return gachaIds, err
}
