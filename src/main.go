package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = InititaliseLogger()

func main() {

	settings := InitialiseSettings()
	eventsStore := NewEventsStore(settings.DatabaseHost, settings.DatabasePort, settings.DatabaseUsername, settings.DatabasePassword, settings.DatabaseName, settings.PageSize)
	handlersManager := NewHandlersManager()
	upgrader := websocket.Upgrader{}
	router := mux.NewRouter()

	InitialiseEventsController(router, eventsStore, handlersManager, settings.PageSize, log)
	InitSubscribersController(router, eventsStore, upgrader, handlersManager, log)

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Starting server at", listenAddr)

	logFatal(fmt.Sprintf("Unable to start server at '%s'", listenAddr), http.ListenAndServe(listenAddr, router))
}
