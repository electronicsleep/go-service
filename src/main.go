package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Usage struct {
	Name  string   `json:"name"`
	Usage []string `json:"usage"`
}

type Health struct {
	Status string `json:"status"`
}

const healthCheckDB = false

var userDatasource = ""
var writerDatasource = ""
var readerDatasource = ""
var datasourcePassword = ""

func init() {
	log.Println("INFO: Init..")

	envVar := "userDatasource"
	userDatasource := os.Getenv(envVar)
	if userDatasource != "" {
		log.Println("INFO: envVar " + envVar + " set")
	} else {
		log.Println("ERROR: envVar " + envVar + " missing")
	}

	envVar = "writerDatasource"
	writerDatasource := os.Getenv(envVar)
	if writerDatasource != "" {
		log.Println("INFO: envVar " + envVar + " set")
	} else {
		log.Println("ERROR: envVar " + envVar + " missing")
	}

	envVar = "readerDatasource"
	readerDatasource := os.Getenv(envVar)
	if readerDatasource != "" {
		log.Println("INFO: envVar " + envVar + " set")
	} else {
		log.Println("ERROR: envVar " + envVar + " missing")
	}

	envVar = "datasourcePassword"
	datasourcePassword := os.Getenv(envVar)
	if datasourcePassword != "" {
		log.Println("INFO: envVar " + envVar + " set")
	} else {
		log.Println("ERROR: envVar " + envVar + " missing")
	}

	log.Println("INFO: Setup DB connection pools")
	OpenDBConn(userDatasource, writerDatasource, datasourcePassword)
	OpenDBRoConn(userDatasource, readerDatasource, datasourcePassword)
	log.Println("INFO: Done init")
}

func checkError(info string, err error) {
	if err != nil {
		fmt.Println("ERROR: ", info)
		log.Println("ERROR: ", err)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	usage := Usage{Name: "go-service", Usage: []string{"DevOps", "SRE", "Infra"}}
	log.Println("INFO: usage: ", usage)

	jsonData, err := json.Marshal(usage)
	if err != nil {
		log.Println("ERROR: Marshal InfoHandler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	// Option to fail health check if DB is not accessable
	if healthCheckDB {
		dbStatusResponse, err := GetEvents("")
		if err != nil {
			log.Println("ERROR: DB Connect issue")
			http.Error(w, dbStatusResponse, http.StatusInternalServerError)
			return
		}
	}

	health := Health{Status: "up"}
	jsonData, err := json.Marshal(health)
	if err != nil {
		log.Println("ERROR: Marshal HealthHandler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("INFO: HealthHandler ", health)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: EventsHandler endpoint")
	query := r.URL.Query()
	service := query.Get("service")
	dbStatusResponse := ""
	var err error
	log.Println("INFO: service:", service)
	if service != "" {
		dbStatusResponse, err = GetEvents(service)
	} else {
		dbStatusResponse, err = GetEvents("")
	}
	if err != nil {
		http.Error(w, dbStatusResponse, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	b := []byte(dbStatusResponse)
	w.Write(b)
}

func eventAddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Endpoint EventAddHandler")
	path := r.URL.Path
	log.Println("INFO: Path:", path)
	if r.Method != "POST" {
		msg := "BadRequest 4xx " + path + "  event requires POST method"
		log.Println("INFO:", msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	type Event struct {
		Event     string
		Service   string
		EventType string
		Datetime  string
	}

	decoder := json.NewDecoder(r.Body)
	var e Event
	err := decoder.Decode(&e)
	log.Println("INFO: Event:", e)
	if err != nil {
		msg := "BadRequest 4xx " + path + " decode error"
		log.Println("INFO:", msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if e.Event == "" {
		msg := "BadRequest 4xx " + path + " event empty"
		log.Println("INFO:", msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else {
		serviceMessage := ""
		dateTimeMessage := ""
		if e.Service == "" {
			e.Service = "None"
		} else {
			serviceMessage = "Service: " + e.Service
		}

		if e.EventType == "" {
			msg := "BadRequest 4xx " + path + " eventType empty"
			log.Println("INFO:", msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		insertResult := ""
		if e.Datetime == "" {
			insertResult, err = InsertEventNow(e.Service, e.Event, e.EventType)
		} else {
			insertResult, err = InsertEvent(e.Service, e.Event, e.EventType, e.Datetime)
		}
		if err != nil {
			http.Error(w, insertResult, http.StatusInternalServerError)
			return
		}
		msg := "AddEvent: " + serviceMessage + " Event: " + e.Event + dateTimeMessage
		log.Println("INFO:", msg)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	healthStatusResponse := "ok"
	log.Println("INFO: StatusHandler")
	fmt.Fprintf(w, "go-service: %s", healthStatusResponse)
}

func main() {
	health := http.HandlerFunc(healthHandler)
	http.Handle("/health", health)

	events := http.HandlerFunc(eventsHandler)
	http.Handle("/events", events)

	eventAdd := http.HandlerFunc(eventAddHandler)
	http.Handle("/add", eventAdd)

	info := http.HandlerFunc(infoHandler)
	http.Handle("/info", info)

	status := http.HandlerFunc(statusHandler)
	http.Handle("/status", status)
	http.Handle("/", status)

	port := "8081"
	log.Println("INFO: Listening...")
	log.Println("INFO: Server: http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
