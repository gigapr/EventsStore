package controllers //import "gigapr/eventsstore/controllers"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gigapr/eventsstore/mappers"
	"gigapr/eventsstore/persistence"
	"gigapr/eventsstore/viewModels"
	"gigapr/eventsstore/websocket"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type eventsController struct {
	EventsStore     *persistence.EventsStore
	HandlersManager *websocket.HandlersManager
	pageSize        int
	log             *logrus.Logger
}

//InitialiseEventsController register all the EventsController routes
func InitialiseEventsController(router *mux.Router, eventStore *persistence.EventsStore, handlersManager *websocket.HandlersManager, pageSize int, log *logrus.Logger) {

	es := new(eventsController)
	es.EventsStore = eventStore
	es.HandlersManager = handlersManager
	es.log = log
	es.pageSize = pageSize

	router.HandleFunc("/event", es.saveEventHandler)
	router.HandleFunc("/{sourceID}/events/{startFrom}/{eventType}", es.getEventsHandler)
}

func (ec *eventsController) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	eventType := vars["eventType"]
	sourceID := vars["sourceID"]
	startFrom, err := strconv.Atoi(vars["startFrom"])
	if err != nil {
		ec.logHttpError(w, r, fmt.Sprintf("'%s' is not a valid number", vars["startFrom"]), http.StatusBadRequest, err)
		return
	}

	eventsStats, err := ec.EventsStore.GetEventsStats(eventType, sourceID)
	if err != nil {
		ec.logHttpError(w, r, fmt.Sprintf("Unable to get events starting from %d.", startFrom), http.StatusInternalServerError, err)
		return
	}

	eventsList := viewModels.EventsResponse{
		Page: viewModels.Page{
			Size:          ec.pageSize,
			TotalElements: eventsStats.Count,
			TotalPages:    (eventsStats.Count + ec.pageSize - 1) / ec.pageSize,
			// Number:        1,
		},
	}

	if startFrom > eventsStats.Count {
		w.WriteHeader(http.StatusNotFound)
	} else {
		eventsList.Links = viewModels.NewLinks(startFrom, eventType, sourceID, eventsStats.Min, eventsStats.Max, eventsStats.Count, ec.pageSize)

		eventsDto, err := ec.EventsStore.Get(startFrom, eventType, sourceID)
		if err != nil {
			ec.logHttpError(w, r, fmt.Sprintf("Unable to get events starting from %d.", startFrom), http.StatusInternalServerError, err)
			return
		}

		eventsList.Embedded = viewModels.Embedded{
			EventsList: mappers.MapEvents(eventsDto),
		}
		w.WriteHeader(http.StatusOK)
	}

	json, err := json.Marshal(eventsList)
	if err != nil {
		ec.logHttpError(w, r, "Unable to encode events response to JSON for subscribers.", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (ec *eventsController) saveEventHandler(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Please send a request body.", http.StatusBadRequest)
		return
	}

	var evm viewModels.Event
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&evm)
	validate := validator.New()

	err = validate.Struct(evm)
	if err != nil {
		var buffer bytes.Buffer

		for _, err := range err.(validator.ValidationErrors) {
			buffer.WriteString(err.Namespace() + " " + err.Tag())
		}

		http.Error(w, buffer.String(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventDataJSON, err := json.Marshal(evm.Data)

	if err != nil {
		ec.logHttpError(w, r, "Unable to encode event data to JSON.", http.StatusInternalServerError, err)
		return
	}

	eventMetadataJSON, err := json.Marshal(evm.Metadata)
	if err != nil {
		ec.logHttpError(w, r, "Unable to encode event metadata to JSON.", http.StatusInternalServerError, err)
		return
	}

	alreadyExist, err := ec.EventsStore.Exists(evm.SourceID, evm.EventID)

	if err != nil {
		ec.logHttpError(w, r, "Unable to persist event.", http.StatusInternalServerError, err)
		return
	}

	if alreadyExist {
		w.WriteHeader(http.StatusOK)
		return
	}

	id, err := ec.EventsStore.Save(evm.SourceID, evm.EventID, evm.Type, eventDataJSON, eventMetadataJSON)

	if err != nil {
		ec.logHttpError(w, r, "Unable to persist event.", http.StatusInternalServerError, err)
		return
	}

	subscribers := ec.HandlersManager.GetChannels(evm.Type)

	if subscribers != nil {
		savedEvent := viewModels.SavedEventResponse{Event: evm}
		savedEvent.Sequence = id

		event, err := json.Marshal(savedEvent)
		if err != nil {
			ec.logHttpError(w, r, "Unable to encode event to JSON for subscribers.", http.StatusInternalServerError, err)
			return
		}

		ec.dispatchToSubscribers(event, subscribers)
	}

	w.WriteHeader(http.StatusCreated)
}

func (ec *eventsController) dispatchToSubscribers(response []byte, handlers []chan []byte) {
	go func(h []chan []byte) {
		for i := range h {
			channel := h[i]
			channel <- response
		}
	}(handlers)
}

func (ec *eventsController) logHttpError(w http.ResponseWriter, r *http.Request, errorMessage string, httpStatusCode int, err error) {
	log := ec.log.WithFields(logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"message":         errorMessage,
	})

	log.Error(r, errorMessage, err)

	http.Error(w, errorMessage, httpStatusCode)
}
