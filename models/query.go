package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func SelectQuery(stringQuery string, args ...interface{}) ([]map[string]string, error) {
	var valueList []map[string]string
	stmt, err := db.Prepare(stringQuery)
	if err != nil {
		return valueList, err
	}
	rows, err := stmt.Query(args...)
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
	return valueList, nil
}

func InsertQuery(stringQuery string, args ...interface{}) (int64, error) {
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
	return id, nil
}

func UpdateQuery(stringQuery string, args ...interface{}) (error) {
	tx,_ := db.Begin()
	stmtIns, err := db.Prepare(stringQuery)
	if err != nil {
		return err
	}
	stmtIns.Exec(args...)
	tx.Commit()
	return nil
}

func DeleteQuery(stringQuery string, args ...interface{}) (error) {
	tx,_ := db.Begin()
	stmtIns, err := db.Prepare(stringQuery)
	if err != nil {
		return err
	}
	stmtIns.Exec(args...)
	tx.Commit()
	return nil
}
