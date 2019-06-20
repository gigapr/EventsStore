package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type subscribersController struct {
	Upgrader        websocket.Upgrader
	HandlersManager *HandlersManager
}

func InitSubscribersController(router *mux.Router, eventStore *EventsStore, upgrader websocket.Upgrader, handlersManager *HandlersManager) {

	sc := new(subscribersController)
	sc.Upgrader = upgrader
	sc.HandlersManager = handlersManager

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
		LogHttpError(w, r, "Unable to serialise list of subscribers to json.", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//subscribe?topic=eventType
func (sc *subscribersController) subscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		LogHttpError(w, r, "Need to specify a topic when subscribing to events.", http.StatusBadRequest, nil)
		return
	}

	c, err := sc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		LogHttpError(w, r, "Unable to upgrade the HTTP server connection to the WebSocket protocol.", http.StatusInternalServerError, err)
		return
	}
	defer c.Close()

	channel := sc.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			LogHttpError(w, r, "Unable to write message to the websocket.", http.StatusInternalServerError, err)
		}
	}
}
