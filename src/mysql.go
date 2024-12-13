package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type EventsTable struct {
	EventId   string `json:"event_id"`
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

func checkErr(err error) {
	if err != nil {
		errStr := "ERROR: Issue connecting to the database"
		log.Printf("%s: %v\n", errStr, err)
	}
}

func OpenDBConn(userDatasource string, writerDatasource string, datasourcePassword string) bool {
	dataSource := userDatasource + ":" + datasourcePassword + "@tcp(" + writerDatasource + ":3306)/infradb"
	db, dbConnErr = sql.Open("mysql", dataSource)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if dbConnErr != nil {
		errStr := "ERROR: DB open issue"
		log.Println(errStr)
		log.Println(dbConnErr)
		return false
	} else {
		err := db.Ping()
		if err != nil {
			log.Println("ERROR: db.Ping failed:", err)
		}
		log.Println("INFO: Opened DB Connection")
		return true
	}
}

func OpenDBRoConn(userDatasource string, readerDatasource string, datasourcePassword string) bool {
	dataSource := userDatasource + ":" + datasourcePassword + "@tcp(" + readerDatasource + ":3306)/infradb"
	dbRo, dbConnErr = sql.Open("mysql", dataSource)
	dbRo.SetMaxOpenConns(maxOpenConns)
	dbRo.SetMaxIdleConns(maxIdleConns)
	dbRo.SetConnMaxLifetime(connMaxLifetime)

	if dbConnErr != nil {
		errStr := "ERROR: DBRo open issue"
		log.Println(errStr)
		log.Println(dbConnErr)
		return false
	} else {
		err := dbRo.Ping()
		if err != nil {
			log.Println("ERROR: dbRo.Ping failed:", err)
		}
		log.Println("INFO: Opened DBRo Connection")
		return true
	}
}

func GetAllEvents() (string, error) {
	log.Println("INFO: GetAllEvents")
	var errStr = ""

	results, err := dbRo.Query("SELECT * FROM events ORDER BY datetime DESC LIMIT 100")
	if err != nil {
		errStr = "ERROR: GetAllEvents: DB Select events issue"
		log.Printf("%s: %v\n", errStr, err)
		return errStr, err
	}
	defer results.Close()
	numEvents := 0
	var eventsList []EventsTable
	var events EventsTable
	for results.Next() {
		err = results.Scan(&events.EventId, &events.Service, &events.Event, &events.EventType, &events.Datetime)
		if err != nil {
			errStr = "ERROR: DB Select events result issue"
			log.Printf("%s: %v\n", errStr, err)
			return errStr, nil
		}
		eventsList = append(eventsList, events)
		numEvents += 1
	}
	jsonData, err := json.Marshal(eventsList)
	if err != nil {
		return "ERROR: GetAllEvents: json.Marshal", err
	}
	return string(jsonData), nil
}

func checkCount(rows *sql.Rows) int {
	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		checkErr(err)
		return count
	}
	return count
}

func InsertEvent(service string, event string, eventType string, datetime string) (string, error) {
	log.Println("INFO: InsertEvent")
	var errStr = ""
	query := "SELECT COUNT(*) as count FROM events WHERE service = ? and event = ? and event_type = ? and datetime = ?"
	results, err := db.Query(query, service, event, eventType, datetime)
	log.Println("INFO:", results)
	numRows := 0
	if err != nil {
		return "ERROR: count events issue", err
	} else {
		numRows = checkCount(results)
	}
	log.Println("INFO: numRows:", numRows)
	if numRows != 0 {
		log.Println("INFO: found duplicate not inserting")
		return "duplicate", nil
	}

	query = "INSERT INTO events (event_id, service, event, event_type, datetime) values (UUID(), ?, ?, ?, ?)"
	result, err := db.Exec(query, service, event, eventType, datetime)

	if err != nil {
		errStr = "ERROR: InsertEvent: DB Insert events issue"
		log.Printf("%s: %v\n", errStr, err)
		return errStr, err
	}
	print(result)
	log.Println("INFO: result: ", result)
	return "ok", nil
}

func InsertEventNow(service string, event string, eventType string) (string, error) {
	log.Println("INFO: InsertEventNow")
	var errStr = ""
	query := "INSERT INTO events (event_id, service, event, event_type, datetime) values (UUID(), ?, ?, ?, NOW())"
	result, err := db.Exec(query, service, event, eventType)
	if err != nil {
		errStr = "ERROR: InsertEventsNow: DB Insert events issue"
		log.Printf("%s: %v\n", errStr, err)
		return errStr, err
	}
	log.Println("INFO: result: ", result)
	return "ok", nil
}
