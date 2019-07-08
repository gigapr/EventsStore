package persistence //import "gigapr/eventsstore/persistence"

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type EventDto struct {
	ID        int
	SourceID  string
	EventID   string
	EventType string
	EventData string
	Metadata  string
	Received  time.Time
}

type eventsStats struct {
	Min   int
	Max   int
	Count int
}

//EventsStore is responsible for storing and retrieving events
type EventsStore struct {
	db       *sql.DB
	pageSize int
}

//NewEventsStore creates an instance of the EventsStore
func NewEventsStore(host string, port int, username string, password string, databaseName string, pageSize int) *EventsStore {

	if pageSize < 1 {
		// log.WithFields(logrus.Fields{"message": "pageSize should be greater than 0."}).Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		// log.WithFields(logrus.Fields{"message": "Unable to open connection to database."}).Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		// log.WithFields(logrus.Fields{"message": "Failed to execute ping against database."}).Fatal(err)
	}

	return &EventsStore{
		db:       db,
		pageSize: pageSize,
	}
}

func (es EventsStore) GetEventsStats(eventType string, sourceID string) (*eventsStats, error) {
	sqlStatement := `select min(id), max(id), count(id)
					 FROM events 
				  	 WHERE LOWER(EventType) = LOWER($1)
				  	 AND LOWER(SourceId) = LOWER($2)`

	rows, err := es.db.Query(sqlStatement, eventType, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var count, min, max int

	for rows.Next() {
		if err := rows.Scan(&min, &max, &count); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &eventsStats{
		Min:   min,
		Max:   max,
		Count: count,
	}, nil
}

func (es EventsStore) Get(sequenceNumber int, eventType string, sourceID string) ([]EventDto, error) {
	events := []EventDto{}
	sqlStatement := `SELECT Id, SourceId, EventId, EventType, EventData, Metadata, Received
					 FROM events
					 WHERE LOWER(EventType) = LOWER($1) 
					 AND LOWER(SourceId) = LOWER($2)
					 AND id > $3
					 LIMIT $4`

	rows, err := es.db.Query(sqlStatement, eventType, sourceID, sequenceNumber-1, es.pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var evn EventDto
		if err := rows.Scan(&evn.ID, &evn.SourceID, &evn.EventID, &evn.EventType, &evn.EventData, &evn.Metadata, &evn.Received); err != nil {
			return nil, err
		}
		events = append(events, evn)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
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
