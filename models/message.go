package models

import (
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Message struct{
	gorm.Model
	content string `gorm:"type:longtext;"`
	publisher User
	create_date time.Time
	read_users []*User `gorm:"many2many:user_messages;"`
	room Room
}
