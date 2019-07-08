package controllers //import "gigapr/eventsstore/controllers"

import (
	"encoding/json"
	gws "gigapr/eventsstore/websocket"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type SubscribersController struct {
	Upgrader        websocket.Upgrader
	HandlersManager *gws.HandlersManager
	Log             *logrus.Logger
}

func (sc *SubscribersController) GetSubscibers(w http.ResponseWriter, r *http.Request) {
	subscribers := sc.HandlersManager.GetAllChannels()

	info := make(map[string]int)

	for k, v := range subscribers {
		info[k] = len(v)
	}

	json, err := json.Marshal(info)

	if err != nil {
		httpError(sc.Log, w, r, "Unable to serialise list of subscribers to json.", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (sc *SubscribersController) Subscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		httpError(sc.Log, w, r, "Need to specify a topic when subscribing to events.", http.StatusBadRequest, nil)
		return
	}

	c, err := sc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		httpError(sc.Log, w, r, "Unable to upgrade the HTTP server connection to the WebSocket protocol.", http.StatusInternalServerError, err)
		return
	}
	defer c.Close()

	channel := sc.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			httpError(sc.Log, w, r, "Unable to write message to the websocket.", http.StatusInternalServerError, err)
		}
	}
}
