package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
)

type Room struct{
	gorm.Model
	slug string `gorm:"NOT NULL;size:255"`
	desc sql.NullString `gorm:"size:255"`
	users []*User `gorm:"many2many:user_rooms;"`
	messages []Message `gorm:"foreignkey:room"`
	is_active bool
}
