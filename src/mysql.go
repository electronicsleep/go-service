package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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
			log.Println("ERROR: DBRo Ping failed:", err)
		}
		log.Println("INFO: Opened DBRo Connection")
		return true
	}
}

func GetEvents(service string) (string, error) {
	log.Println("INFO: GetEvents")
	var errStr = ""
	var err error
	var results *sql.Rows
	if service == "" {
		results, err = dbRo.Query("SELECT * FROM events ORDER BY datetime DESC LIMIT 100")
	} else {
		query := "SELECT * FROM events WHERE service = ? ORDER BY datetime DESC LIMIT 100"
		results, err = dbRo.Query(query, service)
	}
	if err != nil {
		errStr = "ERROR: GetEvents: DB Select events issue"
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
			errStr = "ERROR: GetEvents: DB Select events result issue"
			log.Printf("%s: %v\n", errStr, err)
			return errStr, nil
		}
		eventsList = append(eventsList, events)
		numEvents += 1
	}
	jsonData, err := json.Marshal(eventsList)
	if err != nil {
		return "ERROR: GetEvents: json.Marshal", err
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
	checkQuery := "SELECT COUNT(*) as count FROM events WHERE service = ? and event = ? and event_type = ? and datetime = ?"
	checkResults, err := db.Query(checkQuery, service, event, eventType, datetime)
	numRows := 0
	if err != nil {
		return "ERROR: count events issue", err
	} else {
		numRows = checkCount(checkResults)
	}
	defer checkResults.Close()

	log.Println("INFO: NumRows:", numRows)
	if numRows != 0 {
		log.Println("INFO: FoundDuplicate")
		return "duplicate", nil
	}

	uuid := uuid.New()

	query := "INSERT INTO events (event_id, service, event, event_type, datetime) values (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, uuid, service, event, eventType, datetime)

	if err != nil {
		errStr = "ERROR: InsertEvent: DB Insert events issue"
		log.Printf("%s: %v\n", errStr, err)
		return errStr, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		errStr = "ERROR: InsertEventsNow: DB last insert id issue"
		log.Printf("%s: %v\n", errStr, err)
	}
	log.Println("INFO: insertResult: Rows:", rows)
	return "ok", nil
}

func InsertEventNow(service string, event string, eventType string) (string, error) {
	log.Println("INFO: InsertEventNow")

	var errStr = ""
	checkQuery := "SELECT COUNT(*) as count FROM events WHERE service = ? and event = ? and event_type = ? and datetime = NOW()"
	checkResults, err := db.Query(checkQuery, service, event, eventType)
	numRows := 0
	if err != nil {
		return "ERROR: count events issue", err
	} else {
		numRows = checkCount(checkResults)
	}
	defer checkResults.Close()

	log.Println("INFO: NumRows:", numRows)
	if numRows != 0 {
		log.Println("INFO: FoundDuplicate")
		return "duplicate", nil
	}

	uuid := uuid.New()

	t := time.Now()
	tf := t.Format("2006/01/02 15:04:05")
	log.Println("INFO: TIME:", tf)

	query := "INSERT INTO events (event_id, service, event, event_type, datetime) values (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, uuid, service, event, eventType, tf)
	if err != nil {
		errStr = "ERROR: InsertEventsNow: DB Insert events issue"
		log.Printf("%s: %v\n", errStr, err)
		return errStr, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		errStr = "ERROR: InsertEventsNow: DB last insert id issue"
		log.Printf("%s: %v\n", errStr, err)
	}

	log.Println("INFO: InsertResult: Rows:", rows)
	return "ok", nil
}
