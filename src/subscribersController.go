package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type subscribersController struct {
	log             *HTTPRequestLogger
	Upgrader        websocket.Upgrader
	HandlersManager *HandlersManager
}

func InitSubscribersController(router *mux.Router, eventStore *EventsStore, upgrader websocket.Upgrader, handlersManager *HandlersManager) {

	sc := new(subscribersController)
	sc.Upgrader = upgrader
	sc.HandlersManager = handlersManager
	sc.log = NewHTTPRequestLogger(log)

	router.HandleFunc("/subscribe", sc.subscribe)
	router.HandleFunc("/subscribers", sc.getSubscibers)
}

func (sc *subscribersController) getSubscibers(w http.ResponseWriter, r *http.Request) {
	subscribers := sc.HandlersManager.GetAllChannels()

	info := make(map[string]int)

	for k, v := range subscribers {
		info[k] = len(v)
	}

	json, err := json.Marshal(info)

	if err != nil {
		sc.log.Error(r, "Unable to serialise subscribers to json.", err)
		http.Error(w, "Unable to get the list of subscribers.", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//subscribe?topic=eventType
func (sc *subscribersController) subscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		message := "Need to specify a topic when subscribing to events."
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	c, err := sc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		message := "Unable to upgrade the HTTP server connection to the WebSocket protocol."
		sc.log.Error(r, message, err)
		http.Error(w, message, http.StatusBadRequest)

		return
	}
	defer c.Close()

	channel := sc.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			message := "Unable to write message to the websocket."
			sc.log.Debug(r, "Unsubscribing broken channel.", err)
			go sc.HandlersManager.Unsubscribe(topic, channel)
			http.Error(w, message, http.StatusBadRequest)
		}
	}
}
