package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

func dbInit() (*sql.DB, error){
	db, err := sql.Open(
		"mysql",
		DB_USERNAME+":"+DB_PASSWORD+"@tcp("+DB_HOST+":"+DB_PORT+")/"+DB_NAME+"?charset=utf8")
	if err != nil {
		return db, err
	}
	return db, nil
}


func (m *Models) SelectQuery(stringQuery string, args ...interface{}) ([]map[string]string, error) {
	var valueList []map[string]string
	db, err := dbInit()
	if err != nil {
		return valueList, err
	}
	rows, err := db.Query(stringQuery, args...)
	if err != nil {
		return valueList, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return valueList, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return valueList, err
		}
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
		return valueList, err
	}
	db.Close()
	return valueList, nil
}

func (m *Models) InsertQuery(stringQuery string, args ...interface{}) (int64, error) {
	db, err := dbInit()
	if err != nil {
		return 0,err
	}
	tx,_ := db.Begin()
	stmtIns, err := db.Prepare(stringQuery)
	if err != nil {
		return 0,err
	}
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0,err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0,err
	}
	tx.Commit()
	db.Close()
	return id, nil
}

func (m *Models) UpdateQuery(stringQuery string, args ...interface{}) (error) {
	db, err := dbInit()
	if err != nil {
		return err
	}
	tx,_ := db.Begin()
	stmtIns, err := db.Prepare(stringQuery)
	if err != nil {
		return err
	}
	stmtIns.Exec(args...)
	tx.Commit()
	db.Close()
	return nil
}

func (m *Models) DeleteQuery(stringQuery string, args ...interface{}) (error) {
	db, err := dbInit()
	if err != nil {
		return err
	}
	tx,_ := db.Begin()
	stmtIns, err := db.Prepare(stringQuery)
	if err != nil {
		return err
	}
	stmtIns.Exec(args...)
	tx.Commit()
	db.Close()
	return nil
}
