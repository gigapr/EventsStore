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
	EventsStore *EventsStore
	Upgrader    websocket.Upgrader
	Handler     chan []byte
}

func (ec EventsController) RegisterRoutes() {
	http.HandleFunc("/subscribe", ec.subscribe)
	http.HandleFunc("/event", ec.saveEventHandler)
}

// https://flaviocopes.com/golang-event-listeners/

//subscribe?to=eventType
func (ec EventsController) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := ec.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		message := "unable to upgrades the HTTP server connection to the WebSocket protocol"
		http.Error(w, message, http.StatusBadRequest)
		log.Print(message, err)
		return
	}
	defer c.Close()

	for {
		msg := <-ec.Handler
		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			message := "Unable to write messagge to the websocket"
			log.Println(message, err)
			http.Error(w, message, http.StatusBadRequest)
			break
		}
	}
}

func (ec EventsController) saveEventHandler(w http.ResponseWriter, r *http.Request) {
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

	ec.dispatchToSubscribers(json)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ec EventsController) dispatchToSubscribers(response []byte) {
	go func(handler chan []byte) {
		handler <- response
	}(ec.Handler)
}
