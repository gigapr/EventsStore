package main

import (
	"database/sql"
	"fmt"
	"gigapr/eventsstore/controllers"
	"gigapr/eventsstore/persistence"
	gws "gigapr/eventsstore/websocket"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

//InitSubscribersController register all the EventsController routes
func InitSubscribersController(router *mux.Router, eventStore *persistence.EventsStore, upgrader websocket.Upgrader, handlersManager *gws.HandlersManager, log *logrus.Logger) {

	sc := new(controllers.SubscribersController)
	sc.Upgrader = upgrader
	sc.HandlersManager = handlersManager
	sc.Log = log

	router.HandleFunc("/subscribe", sc.Subscribe)
	router.HandleFunc("/subscribers", sc.GetSubscibers)
}

//InitialiseEventsController register all the EventsController routes
func InitialiseEventsController(router *mux.Router, eventStore *persistence.EventsStore, handlersManager *gws.HandlersManager, pageSize int, log *logrus.Logger) {

	es := new(controllers.EventsController)
	es.EventsStore = eventStore
	es.HandlersManager = handlersManager
	es.Log = log
	es.PageSize = pageSize

	router.HandleFunc("/event", es.SaveEventHandler)
	router.HandleFunc("/{sourceID}/events/{startFrom}/{eventType}", es.GetEventsHandler)
}

//NewHandlersManager inititalise a new HandlersManager
func NewHandlersManager() *gws.HandlersManager {
	es := new(gws.HandlersManager)
	es.Handlers = make(map[string][]chan []byte)
	return es
}

//NewEventsStore creates an instance of the EventsStore
func NewEventsStore(host string, port int, username string, password string, databaseName string, pageSize int) *persistence.EventsStore {

	if pageSize < 1 {
		logFatal("pageSize should be greater than 0.")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, databaseName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logFatal("Unable to open connection to database.", err)
	}

	err = db.Ping()
	if err != nil {
		logFatal("Failed to execute ping against database.", err)
	}

	return &persistence.EventsStore{
		Db:       db,
		PageSize: pageSize,
	}
}

//InititaliseLogger returns an instance of the Logger
func InititaliseLogger() *logrus.Logger {
	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	return log
}
