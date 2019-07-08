package main

import (
	"gigapr/eventsstore/controllers"
	"gigapr/eventsstore/persistence"
	gws "gigapr/eventsstore/websocket"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var log = inititaliseLogger()

func main() {

	settings := InitialiseSettings()
	eventsStore := persistence.NewEventsStore(settings.DatabaseHost, settings.DatabasePort, settings.DatabaseUsername, settings.DatabasePassword, settings.DatabaseName, settings.PageSize)
	handlersManager := gws.NewHandlersManager()
	upgrader := websocket.Upgrader{}
	router := mux.NewRouter()

	controllers.InitialiseEventsController(router, eventsStore, handlersManager, settings.PageSize, log)
	controllers.InitSubscribersController(router, eventsStore, upgrader, handlersManager, log)

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, router))
}

func inititaliseLogger() *logrus.Logger {
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
