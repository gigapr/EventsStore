package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type eventViewModel struct {
	SourceID string      `json:"sourceId"`
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
}

type EventsController struct {
	log             *HTTPRequestLogger
	EventsStore     *EventsStore
	Upgrader        websocket.Upgrader
	HandlersManager *HandlersManager
}

func NewEventsController(eventStore *EventsStore, upgrader websocket.Upgrader, handlersManager *HandlersManager) *EventsController {

	es := new(EventsController)
	es.EventsStore = eventStore
	es.Upgrader = upgrader
	es.HandlersManager = handlersManager
	es.log = NewHTTPRequestLogger()

	es.registerRoutes()

	return es
}

func (ec *EventsController) registerRoutes() {
	http.HandleFunc("/subscribe", ec.subscribe)
	http.HandleFunc("/event", ec.saveEventHandler)
	http.HandleFunc("/subscribers", ec.getSubscibers)
}

func (ec *EventsController) getSubscibers(w http.ResponseWriter, r *http.Request) {
	subscribers := ec.HandlersManager.GetAllChannels()

	info := make(map[string]int)

	for k, v := range subscribers {
		info[k] = len(v)
	}

	json, err := json.Marshal(info)

	if err != nil {
		ec.log.Error(r, "Unable to serialise subscribers to json.", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//subscribe?topic=eventType
func (ec *EventsController) subscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		message := "Need to specify a topic when subscribing to events."
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	c, err := ec.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		message := "Unable to upgrade the HTTP server connection to the WebSocket protocol."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusBadRequest)

		return
	}
	defer c.Close()

	channel := ec.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			message := "Unable to write message to the websocket."
			ec.log.Debug(r, "Unsubscribing broken channel.", err)
			go ec.HandlersManager.Unsubscribe(topic, channel)
			http.Error(w, message, http.StatusBadRequest)
		}
	}
}

func (ec *EventsController) saveEventHandler(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Please send a request body.", http.StatusBadRequest)
		return
	}

	var evm eventViewModel
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&evm)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(evm.Data)
	if err != nil {
		message := "Unable to encode event to JSON."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)

		return
	}

	err = ec.EventsStore.Save(evm.SourceID, evm.Type, json)

	if err != nil {
		message := "Unable to persist event."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)

		return
	}
	subscribers := ec.HandlersManager.GetChannels(evm.Type)

	if subscribers != nil {
		ec.dispatchToSubscribers(json, subscribers)
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
