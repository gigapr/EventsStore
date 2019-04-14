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

func (ec EventsController) subscribe(w http.ResponseWriter, r *http.Request) {
	c, err := ec.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		msg := <-ec.Handler
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (ec EventsController) emit(response []byte) {
	go func(handler chan []byte) {
		handler <- response
	}(ec.Handler)
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

	ec.emit(json)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
