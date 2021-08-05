package models

func (user *UserAuth) Get() error {
	return db.First(&user).Error
}

func (user *UserAuth) Create() error {
	return db.Create(user).Error
}

func (user *UserAuth) Patch() error {
	return db.Save(&user).Error
}

func (user *UserAuth) Delete() error {
	return db.Delete(&user).Error
}
