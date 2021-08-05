package models

import (
	"time"

	"gorm.io/gorm"

	pq "github.com/lib/pq"

	"github.com/tnyie/journaler-api/database"
)

var db *gorm.DB

// TODO Ensure username has no '@' symbol
type User struct {
	ID       string `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"id,omitempty"`
	Username string `gorm:"unique" json:"username,omitempty"`
	Email    string `gorm:"unique" json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	External bool   `json:"external,omitempty"`
}

type UserAuth struct {
	ID       string `gorm:"primaryKey;type:string;default:uuid_generate_v4()" json:"id,omitempty"`
	Email    string `gorm:"unique" json:"email,omitempty"`
	Username string `gorm:"unique" json:"username,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Hash     []byte `json:"hash,omitempty"`
}

type Login struct {
	UserAuth
	Password string `json:"password"`
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

type UserSession struct {
	ID      string
	Expires time.Time
}

// InitModels migrates modesls and initiates database connection
func InitModels() {
	db = database.InitDB()

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";") // enable uuid generation on server
	db.AutoMigrate(&User{}, &Journal{}, &Entry{})
}
