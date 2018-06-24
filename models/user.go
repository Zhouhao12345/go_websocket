package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
)

type User struct{
	gorm.Model
	username string `gorm:"NOT NULL;size:255"`
	password string `gorm:"NOT NULL;size:255"`
	mobile sql.NullString `gorm:"size:255"`
	avator_image sql.NullString `gorm:"size:255"`
	rooms []*Room `gorm:"many2many:user_rooms;"`
	messages []Message `gorm:"foreignkey:publisher"`
	read_messages []*Message `gorm:"many2many:user_messages;"`
	is_active bool
}
