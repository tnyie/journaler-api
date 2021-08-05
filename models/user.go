package models

import (
	"strings"

	goaway "github.com/TwinProduction/go-away"
)

func (user *User) Sanitize() (bool, string) {
	if strings.Contains(user.Username, "@") {
		return false, "invalid username"
	}

	replacer := strings.NewReplacer(" ", "")

	username := replacer.Replace(user.Username)

	user.Username = username

	if goaway.IsProfane(user.Username) {
		return false, "profane username"
	}

	return true, ""
}

func (user *User) Get() error {
	return db.First(&user).Error
}

func (user *User) Create() error {
	err := db.Create(&user).Error

	defaultJournal := &Journal{
		OwnerID:  user.ID,
		ParentID: "",
		Name:     "_default",
	}

	defaultJournal.Create()
	return err
}

func (user *User) Patch() error {
	return db.Save(&user).Error
}

func (user *User) Delete() error {
	return db.Delete(&user).Error
}
