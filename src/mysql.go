package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type EventsTable struct {
	EventId   int    `json:"event_id"`
	Service   string `json:"service"`
	Event     string `json:"event"`
	EventType string `json:"event_type"`
	Datetime  string `json:"datetime"`
}

var maxOpenConns = 20
var maxIdleConns = 5
var connMaxLifetime = time.Minute
var db *sql.DB
var dbRo *sql.DB
var dbConnErr error

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
	}
	return count
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func OpenDBConn(writerDatasource string, datasourcePassword string) bool {
	dataSource := "infradb:" + datasourcePassword + "@tcp(" + writerDatasource + ":3306)/infradb"
	db, dbConnErr = sql.Open("mysql", dataSource)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if dbConnErr != nil {
		errStr := "ERROR: DB connect issue"
		log.Println(errStr)
		log.Println(dbConnErr)
		return false
	} else {
		log.Println("Opened DB Connection")
		return true
	}
}

func OpenDBRoConn(readerDatasource string, datasourcePassword string) bool {
	dataSource := "infradb:" + datasourcePassword + "@tcp(" + readerDatasource + ":3306)/infradb"
	dbRo, dbConnErr = sql.Open("mysql", dataSource)
	dbRo.SetMaxOpenConns(maxOpenConns)
	dbRo.SetMaxIdleConns(maxIdleConns)
	dbRo.SetConnMaxLifetime(connMaxLifetime)

	if dbConnErr != nil {
		errStr := "ERROR: DB Ro connect issue"
		log.Println(errStr)
		log.Println(dbConnErr)
		return false
	} else {
		log.Println("Opened DB Ro Connection")
		return true
	}
}

func GetAllEvents() string {
	log.Println("INFO: GetAllEvents")
	var errStr = ""

	results, err := dbRo.Query("SELECT * FROM events order by datetime desc")
	if err != nil {
		errStr = "ERROR: DB Select events issue"
		log.Println(errStr)
		log.Println(err)
		return errStr
	}
	defer results.Close()
	columns, err := results.Columns()
	if err != nil {
		return "Error: columns"
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for results.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		results.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "Error: json.Marshal"
	}
	return string(jsonData)
}

func InsertEvent(service string, event string, eventType string, datetime string) string {
	log.Println("INFO: InsertEvent")
	var errStr = ""
	query := "SELECT COUNT(*) as count FROM events WHERE service = ? and event = ? and event_type = ? and datetime = ?"
	results, err := db.Query(query, service, event, eventType, datetime)
	log.Println("INFO:", results)
	numRows := checkCount(results)
	log.Println("INFO: numRows:", numRows)
	if numRows != 0 {
		log.Println("INFO: found duplicate not inserting")
		return "duplicate"
	}

	query = "INSERT INTO events (service, event, event_type, datetime) values (?, ?, ?, ?)"
	result, err := db.Exec(query, service, event, eventType, datetime)

	if err != nil {
		errStr = "ERROR: DB Insert events issue"
		log.Println(errStr)
		log.Println(err)
		return errStr
	}
	print(result)
	log.Println("INFO: result: ", result)
	return "ok"
}

func InsertEventNow(service string, event string, eventType string) string {
	log.Println("INFO: InsertEventNow")
	var errStr = ""
	query := "INSERT INTO events (service, event, event_type, datetime) values (?, ?, ?, NOW())"
	result, err := db.Exec(query, service, event, eventType)
	if err != nil {
		errStr = "ERROR: DB Insert events issue"
		log.Println(errStr)
		log.Println(err)
		return errStr
	}
	log.Println("INFO: result: ", result)
	return "ok"
}
