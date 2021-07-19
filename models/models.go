package models

import (
	"gorm.io/gorm"

	pq "github.com/lib/pq"

	"github.com/tnyie/journaler-api/database"
)

var db *gorm.DB

type User struct {
	ID      string `json:"id,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// Journal structure to hold child journals and entries
type Journal struct {
	ID       string         `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"id,omitempty"`
	OwnerID  string         `json:"owner_id,omitempty"`
	ParentID string         `json:"parent_id,omitempty"`
	Name     string         `json:"name,omitempty"`
	Children pq.StringArray `gorm:"type:text[]" json:"children,omitempty"`
	Entries  pq.StringArray `gorm:"type:text[]" json:"entries,omitempty"`
}

// Entry structure with text content
type Entry struct {
	ID        string `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"id,omitempty"`
	JournalID string `json:"journal_id,omitempty"`
	OwnerID   string `json:"owner_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
}

// InitModels migrates modesls and initiates database connection
func InitModels() {
	db = database.InitDB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";") // enable uuid generation on server
	db.AutoMigrate(&User{}, &Journal{}, &Entry{})
}
