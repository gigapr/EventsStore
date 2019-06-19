package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const pageSize = 5

type eventDto struct {
	ID        int
	SourceID  string
	EventID   string
	EventType string
	EventData string
	Metadata  string
	Received  time.Time
}

//EventsStore is responsible for storing and retrieving events
type EventsStore struct {
	log *logrus.Logger
	db  *sql.DB
}

//NewEventsStore creates an instance of the EventsStore
func NewEventsStore(host string, port int, username string, password string, databaseName string) *EventsStore {

	l := logrus.New()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		l.WithFields(logrus.Fields{"message": "Unable to open connection to database."}).Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		l.WithFields(logrus.Fields{"message": "Failed to execute ping against database."}).Fatal(err)
	}

	return &EventsStore{
		log: l,
		db:  db,
	}
}

func (es EventsStore) Get(sequenceNumber int) []eventDto {
	events := []eventDto{}
	sqlStatement := `SELECT Id, SourceId, EventId, EventType, EventData, Metadata, Received from events 
					 WHERE id > $1
					 LIMIT $2`

	rows, err := es.db.Query(sqlStatement, sequenceNumber, pageSize)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var evn eventDto
		if err := rows.Scan(&evn.ID, &evn.SourceID, &evn.EventID, &evn.EventType, &evn.EventData, &evn.Metadata, &evn.Received); err != nil {
			log.Fatal(err)
		}
		events = append(events, evn)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return events
}

//Save an event to the database
func (es EventsStore) Save(sourceID string, EventID string, eventType string, data []byte, metadata []byte) (int, error) {

	sqlStatement := `INSERT INTO Events (SourceId, EventId, EventType, EventData, Metadata)
					 VALUES 			($1, 	   $2, 		$3, 	   $4, 		  $5)
					 RETURNING id`
	id := 0
	err := es.db.QueryRow(sqlStatement, sourceID, EventID, eventType, data, metadata).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// Exists checks if an event already exists in the store
func (es EventsStore) Exists(sourceID string, EventID string) (bool, error) {

	sqlStatement := `SELECT * FROM Events 
					 WHERE SourceId = $1 AND EventId = $2
					 fetch first 1 rows only`

	rows, err := es.db.Query(sqlStatement, sourceID, EventID)
	defer rows.Close()

	if err != nil {
		return true, err
	}

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
