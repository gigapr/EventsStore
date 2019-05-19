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

// needs an endpoint to unsubescribe

//subscribe?topic=eventType
func (ec *EventsController) subscribe(w http.ResponseWriter, r *http.Request) {

	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		message := "need to specify a topic when subscribing to events"
		http.Error(w, message, http.StatusBadRequest)
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

	channel := ec.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			message := "Unable to write messagge to the websocket"
			log.Println(message, err)
			log.Println("Unsubscribing broken channel")
			go ec.HandlersManager.Unsubscribe(topic, channel)
			http.Error(w, message, http.StatusBadRequest)
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

	subscribers := ec.HandlersManager.GetChannels(evm.Type)

	if subscribers != nil {
		ec.dispatchToSubscribers(json, subscribers)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ec *EventsController) dispatchToSubscribers(response []byte, handlers []chan []byte) {
	go func(h []chan []byte) {
		for i := range h {
			channel := h[i]
			channel <- response
		}
	}(handlers)
}
