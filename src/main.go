package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Profile struct {
	Name    string   `json:"name"`
	Hobbies []string `json:"hobbies"`
}

type Health struct {
	Status string `json:"status"`
}

var writerDatasource = ""
var readerDatasource = ""
var datasourcePassword = ""

func init() {
	log.Println("INFO: Init..")

	envVar := "writerDatasource"
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
	OpenDBConn(writerDatasource, datasourcePassword)
	OpenDBRoConn(readerDatasource, datasourcePassword)
	log.Println("INFO: Done init")
}

func checkError(info string, err error) {
	if err != nil {
		fmt.Println("ERROR: ", info)
		log.Println("ERROR: ", err)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	profile := Profile{Name: "Chris", Hobbies: []string{"Bass", "Programming"}}
	log.Println("INFO: api ", profile)

	js, err := json.Marshal(profile)
	if err != nil {
		log.Println("ERROR: Marshal apiHandler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := Health{Status: "up"}
	log.Println("INFO: status ", health)

	js, err := json.Marshal(health)
	if err != nil {
		log.Println("ERROR: Marshal healthHandler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: eventsHandler endpoint")
	dbStatusResponse := ""
	dbStatusResponse = GetAllEvents()
	w.Header().Set("Content-Type", "application/json")
	b := []byte(dbStatusResponse)
	w.Write(b)
}

func eventAddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("INFO: Endpoint eventAddHandler")
	if r.Method != "POST" {
		message := "Add event requires POST method"
		fmt.Fprintf(w, message)
		log.Println("INFO" + message)
		return
	}

	type Event struct {
		Event     string
		Service   string
		EventType string
		Datetime  string
	}

	path := r.URL.Path
	log.Println("INFO: path: " + path)
	decoder := json.NewDecoder(r.Body)
	var e Event
	err := decoder.Decode(&e)
	log.Println("INFO: event ", e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message := ""

	if e.Event == "" {
		fmt.Fprintf(w, "Event blank")
		log.Println("INFO: Event blank")
	} else {
		serviceMessage := ""
		dateTimeMessage := ""
		if e.Service == "" {
			e.Service = "na"
		} else {
			serviceMessage = "Service: " + e.Service
		}

		if e.EventType == "" {
			fmt.Fprintf(w, "ERROR: type not sent")
			return
		}

		insertResult := ""
		if e.Datetime == "" {
			insertResult = InsertEventNow(e.Service, e.Event, e.EventType)
		} else {
			insertResult = InsertEvent(e.Service, e.Event, e.EventType, e.Datetime)
		}
		fmt.Fprintf(w, "Add Event: %s Result: %s", e.Event, insertResult)
		//fmt.Fprintf(w, "Add Event: #{e.Event} Result: #{insertResult}")
		message = "Add Event: " + serviceMessage + " Event " + e.Event + dateTimeMessage
		log.Println("INFO:", message)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	healthStatusResponse := "ok"
	log.Println("INFO: statusHandler")
	fmt.Fprintf(w, "go-service: %s", healthStatusResponse)
}

func main() {
	health := http.HandlerFunc(healthHandler)
	http.Handle("/health", health)

	events := http.HandlerFunc(eventsHandler)
	http.Handle("/events", events)

	eventAdd := http.HandlerFunc(eventAddHandler)
	http.Handle("/add", eventAdd)

	api := http.HandlerFunc(apiHandler)
	http.Handle("/api", api)

	status := http.HandlerFunc(statusHandler)
	http.Handle("/status", status)
	http.Handle("/", status)

	port := "8080"
	log.Println("INFO: Listening...")
	log.Println("INFO: Server: http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
