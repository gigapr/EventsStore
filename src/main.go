package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	eventsController EventsController
)

func main() {

	settings := InitialiseSettings()

	eventsStore := NewEventsStore(settings.DatabaseHost, settings.DatabasePort, settings.DatabaseUsername, settings.DatabasePassword, settings.DatabaseName)
	handlersManager := NewHandlersManager()
	upgrader := websocket.Upgrader{}

	NewEventsController(eventsStore, upgrader, handlersManager)

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
