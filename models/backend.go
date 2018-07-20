package models

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"go_ws/config"
	"log"
	"time"
)

const (
       DB_NAME = config.GLOBAL_DB_NAME
       DB_USERNAME = config.GLOBAL_DB_USERNAME
       DB_PASSWORD = config.GLOBAL_DB_PASSWORD
       DB_HOST = config.GLOBAL_DB_HOST
       DB_PORT = config.GLOBAL_DB_PORT
       DB_CHARSET = config.GLOBAL_DB_CHARSET
       DB_DIRVER = config.GLOBAL_DB_DIRVER
)

var db = new(sql.DB)

func init() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("error: %v", err)
		}
	}()
	var err error
	db, err = sql.Open(DB_DIRVER,
		DB_USERNAME+":"+DB_PASSWORD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME+"?charset="+DB_CHARSET)
	db.SetConnMaxLifetime(time.Second * 60)
	if err != nil {
		panic(err)
	}
}

func DBClose()  {
	db.Close()
}