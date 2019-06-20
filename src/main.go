package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	log := InititaliseLogger()
	settings := InitialiseSettings()
	eventsStore := NewEventsStore(settings.DatabaseHost, settings.DatabasePort, settings.DatabaseUsername, settings.DatabasePassword, settings.DatabaseName)
	handlersManager := NewHandlersManager()
	upgrader := websocket.Upgrader{}
	router := mux.NewRouter()

	InitialiseEventsController(router, eventsStore, handlersManager)
	InitSubscribersController(router, eventsStore, upgrader, handlersManager)

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, router))
}
