package models

func (journal *Journal) Get() error {
	return db.First(&journal).Error
}

func GetOwnJournals(ownerID string) (*[]Journal, error) {
	var journals []Journal
	err := db.Where("owner_id=?", ownerID).Find(&journals).Error
	return &journals, err
}

func (journal *Journal) Create() error {
	return db.Create(&journal).Error
}

func (journal *Journal) AddChild(id string) error {
	journal.Get()
	journal.Children = append(journal.Children, id)
	return db.Save(&journal).Error
}

func (journal *Journal) AddEntry(id string) error {
	journal.Get()
	journal.Entries = append(journal.Entries, id)
	return db.Save(&journal).Error
}
