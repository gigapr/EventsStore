package main

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

//EventsStore is responsible for storing and retrieving events
type EventsStore struct {
	log          *logrus.Logger
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
		log:          logrus.New(),
	}
}

//Save stores an event to the database
func (es EventsStore) Save(sourceID string, EventID string, eventType string, data []byte, metadata []byte) (int, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		es.host, es.port, es.username, es.password, es.databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		es.log.WithFields(logrus.Fields{
			"message": "Unable to open connection to database.",
		}).Error(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO Events (SourceId, EventId, EventType, EventData, Metadata)
					 VALUES 			($1, 	   $2, 		$3, 	   $4, 		  $5)
					 RETURNING id`
	id := 0
	err = db.QueryRow(sqlStatement, sourceID, EventID, eventType, data, metadata).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
