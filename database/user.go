package database

import (
	"fmt"
	"game-api-gin/model"
)

// usersテーブルにユーザ情報を新規追加
func (d *GormDatabase) CreateUser(user model.User) error {
	//	INSERT INTO `users` (`user_id`,`name`,`private_key`)
	//	VALUES ('95daec2b-287c-4358-ba6f-5c29e1c3cbdf','aaa','6e7eada90afb7e84bf5b4498c6adaa2d4014904644637d5fb355266944fbf93a')
	return d.DB.Create(&user).Error
}

// usersテーブルからユーザIDが引数userIdのユーザの情報を取得
func (d *GormDatabase) GetUser(userId string) (model.User, error) {
	var user model.User
	// SELECT * FROM `users` WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	err := d.DB.Where("user_id = ?", userId).Find(&user).Error
	if err != nil {
		return model.User{}, err
	} else if (user == model.User{}) {
		return model.User{}, fmt.Errorf("no data")
	} else {
		return user, err
	}
}

// usersテーブルからユーザIDが引数userIdのユーザの情報を、引数userのものに更新
func (d *GormDatabase) UpdateUser(user model.User, userId string) error {
	// UPDATE `users` SET `name`='bbb' WHERE user_id = '95daec2b-287c-4358-ba6f-5c29e1c3cbdf'
	err := d.DB.Model(&user).Where("user_id = ?", userId).Update("name", user.Name).Error
	if err != nil {
		return err
	} else if (user == model.User{}) {
		return fmt.Errorf("no data")
	} else {
		return err
	}
}
