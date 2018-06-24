package app

import (
	"github.com/jinzhu/gorm"
	"go_ws/config"
	"log"
	"go_ws/models"
)

func Migrate()  {
	db, err := gorm.Open(
		"mysql",
		config.GLOBAL_DB_USERNAME+ ":"+
			config.GLOBAL_DB_PASSWORD+
			"@tcp("+config.GLOBAL_DB_HOST+
			":"+config.GLOBAL_DB_PORT+
			")/"+config.GLOBAL_DB_NAME+
			"?charset="+config.GLOBAL_DB_CHARSET+
			"&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	tx := db.Begin()
	db.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Room{},
	)
	tx.Commit()
	db.Close()
}