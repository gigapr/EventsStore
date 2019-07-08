package controllers //import "gigapr/eventsstore/controllers"

import (
	"encoding/json"
	"gigapr/eventsstore/persistence"
	gws "gigapr/eventsstore/websocket"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type subscribersController struct {
	Upgrader        websocket.Upgrader
	HandlersManager *gws.HandlersManager
	log             *logrus.Logger
}

func InitSubscribersController(router *mux.Router, eventStore *persistence.EventsStore, upgrader websocket.Upgrader, handlersManager *gws.HandlersManager, log *logrus.Logger) {

	sc := new(subscribersController)
	sc.Upgrader = upgrader
	sc.HandlersManager = handlersManager
	sc.log = log

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
		sc.logHttpError(w, r, "Unable to serialise list of subscribers to json.", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//subscribe?topic=eventType
func (sc *subscribersController) subscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if len(topic) == 0 {
		sc.logHttpError(w, r, "Need to specify a topic when subscribing to events.", http.StatusBadRequest, nil)
		return
	}

	c, err := sc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		sc.logHttpError(w, r, "Unable to upgrade the HTTP server connection to the WebSocket protocol.", http.StatusInternalServerError, err)
		return
	}
	defer c.Close()

	channel := sc.HandlersManager.Subscribe(topic)

	for {
		msg := <-channel

		err = c.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			sc.logHttpError(w, r, "Unable to write message to the websocket.", http.StatusInternalServerError, err)
		}
	}
}

func (sc *subscribersController) logHttpError(w http.ResponseWriter, r *http.Request, errorMessage string, httpStatusCode int, err error) {
	log := sc.log.WithFields(logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"message":         errorMessage,
	})

	log.Error(r, errorMessage, err)

	http.Error(w, errorMessage, httpStatusCode)
}
