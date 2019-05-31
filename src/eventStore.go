package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//EventsStore is responsible for storing and retrieving events
type EventsStore struct {
	host         string
	port         int
	username     string
	password     string
	databaseName string
}

//NewEventsStore creates an instance of the EventsStore
func NewEventsStore(host string, port int, username string, password string, databaseName string) *EventsStore {

	return &EventsStore{
		host:         host,
		port:         port,
		username:     username,
		password:     password,
		databaseName: databaseName,
	}
}

//Save stores an event to the database
func (es EventsStore) Save(sourceID string, eventType string, data []byte) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		es.host, es.port, es.username, es.password, es.databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Print(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO Events (SourceId, EventType, EventData)
					 VALUES ($1, $2, $3)
					 RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, sourceID, eventType, data).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
}
