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

	settings := InitialiseSetting()

	eventsController.EventsStore = NewEventsStore()
	eventsController.Upgrader = websocket.Upgrader{} // use default options
	eventsController.Handler = make(chan []byte)
	eventsController.RegisterRoutes()

	http.Handle("/metrics", promhttp.Handler())

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
