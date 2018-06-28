package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"database/sql"
	"encoding/base64"
	"log"
)

type User struct{
	gorm.Model
	Username *string `gorm:"not null;unique"`
	Password *string `gorm:"not null;unique"`
	Mobile sql.NullString `gorm:"size:255"`
	Avator_image sql.NullString `gorm:"size:255"`
	Rooms []*Room `gorm:"many2many:user_rooms;"`
	Messages []Message `gorm:"foreignkey:Publisher"`
	Read_messages []*Message `gorm:"many2many:user_messages;"`
	Is_active bool
}

func (user *User)Auth(password string, username string, DB *gorm.DB) bool {
	DB.Where("username = ?", username).First(&user)
	password_user := *user.Password
	decodeBytes, err := base64.StdEncoding.DecodeString(password_user)
	if err != nil {
		log.Printf(err.Error())
		return false
	}
	if string(decodeBytes) != password {
		return false
	}
	return true
}

func (user *User)Register(password string, username string, DB *gorm.DB)  {
	password_base64 := base64.StdEncoding.EncodeToString([]byte(password))
	user_create := &User{Password:&password_base64, Username:&username, Is_active:true}
	DB.Create(user_create)
}
