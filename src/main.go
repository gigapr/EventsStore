package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {

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

	settings := InitialiseSettings()
	eventsStore := NewEventsStore(settings.DatabaseHost, settings.DatabasePort, settings.DatabaseUsername, settings.DatabasePassword, settings.DatabaseName)
	handlersManager := NewHandlersManager()
	upgrader := websocket.Upgrader{}

	RegisterEventsControllerRoutes(eventsStore, upgrader, handlersManager)
	RegisterSubscribersControllerRoutes(eventsStore, upgrader, handlersManager)

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
