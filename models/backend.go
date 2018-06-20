package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"go_ws/config"
)

const (
	DB_NAME = config.GLOBAL_DB_NAME
	DB_USERNAME = config.GLOBAL_DB_USERNAME
	DB_PASSWORD = config.GLOBAL_DB_PASSWORD
	DB_HOST = config.GLOBAL_DB_HOST
	DB_PORT = config.GLOBAL_DB_PORT
)

type Models struct {
}

func dbInit() *sql.DB{
	db, err := sql.Open(
		"mysql",
		DB_USERNAME+":"+DB_PASSWORD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME+"?charset=utf8")
	checkErr(err)
	return db
}


func (m *Models) SelectQuery(stringQuery string) []map[string]string {
	db := dbInit()
	rows, err := db.Query(stringQuery)
	checkErr(err)
	columns, err := rows.Columns()
	checkErr(err)
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var valueList []map[string]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		checkErr(err)
		var value string
		row := make(map[string]string)
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			row[columns[i]] = value
		}
		valueList = append(valueList, row)
	}
	if err = rows.Err(); err != nil {
		log.Fatalln(err.Error())
	}
	db.Close()
	return valueList
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
