package models

import (
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Message struct{
	gorm.Model
	Content string `gorm:"type:longtext;"`
	Publisher uint
	Read_users []*User `gorm:"many2many:user_messages;"`
	Room uint
}
