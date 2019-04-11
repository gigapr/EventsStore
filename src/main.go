package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	eventsController EventsController
	eventPublisher   EventPublisher
)

var upgrader = websocket.Upgrader{} // use default options

func subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {

		time.Sleep(2 * time.Second)
		s := "Current Unix Time: %v\n" + time.Now().String()

		err = c.WriteMessage(websocket.TextMessage, []byte(s))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {

	settings := InitialiseSetting()

	eventsController.EventPublisher = NewEventPublisher(settings.BrokerConnectionString)
	eventsController.EventsStore = NewEventsStore()
	eventsController.RegisterRoutes()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/subscribe", subscribe)

	listenAddr := "0.0.0.0:" + settings.Port

	log.Println("Server is ready to handle requests at", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
