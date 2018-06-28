package middleware

import (
	"net/http"
	"go_ws/tools"
	"github.com/jinzhu/gorm"
	"go_ws/config"
	"log"
)

func WithDatabaseInit(next tools.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		if config.ATOMIC_REQUEST {
			tx := db.Begin()
			next.ServeHTTP(w,r)
			tx.Commit()
		} else {
			next.ServeHTTP(w,r)
		}
		db.Close()
	}
}
