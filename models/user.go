package models

func (user *User) Get() error {
	return db.First(&user).Error
}

func (user *User) Enable() error {
	user.Enabled = true
	return db.Save(&user).Error
}

func (user *User) Create() error {
	return db.Create(&user).Error
}
