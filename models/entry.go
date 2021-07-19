package models

func (entry *Entry) Get() error {
	return db.First(&entry).Error
}

func (entry *Entry) Create() error {
	return db.Create(&entry).Error
}
