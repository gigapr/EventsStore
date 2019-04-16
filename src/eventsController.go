package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type EventViewModel struct {
	SourceId string      `json:"sourceId"`
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
}

type EventsController struct {
	EventsStore     *EventsStore
	Upgrader        websocket.Upgrader
	HandlersManager *HandlersManager
}

func NewEventsController(eventStore *EventsStore, upgrader websocket.Upgrader, handlersManager *HandlersManager) *EventsController {

	es := new(EventsController)
	es.EventsStore = eventStore
	es.Upgrader = upgrader
	es.HandlersManager = handlersManager
	es.RegisterRoutes()

	return es
}

func (ec *EventsController) RegisterRoutes() {
	http.HandleFunc("/subscribe", ec.subscribe)
	http.HandleFunc("/event", ec.saveEventHandler)
}

//subscribe?topic=eventType
func (ec *EventsController) subscribe(w http.ResponseWriter, r *http.Request) {

	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		message := "need to specify a topic when subscribing to events"
		http.Error(w, message, http.StatusBadRequest)
		log.Print(message)
		return
	}

	c, err := ec.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		message := "unable to upgrades the HTTP server connection to the WebSocket protocol"
		http.Error(w, message, http.StatusBadRequest)
		log.Print(message, err)
		return
	}
	defer c.Close()

	channel := ec.HandlersManager.Get(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			message := "Unable to write messagge to the websocket"
			log.Println(message, err)
			http.Error(w, message, http.StatusBadRequest)
			break
		}
	}
}

func (ec *EventsController) saveEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	var evm EventViewModel
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&evm)

	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(evm.Data)

	ec.EventsStore.Save(evm.SourceId, evm.Type, json)

	ec.dispatchToSubscribers(json, ec.HandlersManager.Get(evm.Type))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ec *EventsController) dispatchToSubscribers(response []byte, handler chan []byte) {
	go func(h chan []byte) {
		log.Println(h)
		h <- response
	}(handler)
}
