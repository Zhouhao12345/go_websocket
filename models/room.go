package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
)

type Room struct{
	gorm.Model
	Slug *string `gorm:"not null;unique"`
	Desc sql.NullString `gorm:"size:255"`
	Users []*User `gorm:"many2many:user_rooms;"`
	Messages []Message `gorm:"foreignkey:Room"`
	Is_active bool
}
