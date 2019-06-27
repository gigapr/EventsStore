package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
)

type eventsController struct {
	EventsStore     *EventsStore
	HandlersManager *HandlersManager
}

//InitialiseEventsController register all the EventsController routes
func InitialiseEventsController(router *mux.Router, eventStore *EventsStore, handlersManager *HandlersManager) {

	es := new(eventsController)
	es.EventsStore = eventStore
	es.HandlersManager = handlersManager

	router.HandleFunc("/event", es.saveEventHandler)
	router.HandleFunc("/events/{startFrom}/{eventType}", es.getEventsHandler)
}

//events?
func (ec *eventsController) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	startFrom, err := strconv.Atoi(vars["startFrom"])
	if err != nil {
		LogHttpError(w, r, fmt.Sprintf("'%s' is not a valid number", vars["startFrom"]), http.StatusBadRequest, err)
		return
	}
	eventType := vars["eventType"]
	eventsDto, err := ec.EventsStore.Get(startFrom, eventType)
	if err != nil {
		LogHttpError(w, r, fmt.Sprintf("Unable to get events starting from %d.", startFrom), http.StatusInternalServerError, err)
		return
	}

	events := mapEvents(eventsDto)

	count := 0
	if len(eventsDto) > 0 {
		count = eventsDto[0].Count
	}

	links := newLinks(startFrom, eventType, count)
	eventsList := newEventsResponse(events, links, pageSize, count, 1, 1)
	json, err := json.Marshal(eventsList)
	if err != nil {
		LogHttpError(w, r, "Unable to encode events response to JSON for subscribers.", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (ec *eventsController) saveEventHandler(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		http.Error(w, "Please send a request body.", http.StatusBadRequest)
		return
	}

	var evm event
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
		LogHttpError(w, r, "Unable to encode event data to JSON.", http.StatusInternalServerError, err)
		return
	}

	eventMetadataJSON, err := json.Marshal(evm.Metadata)
	if err != nil {
		LogHttpError(w, r, "Unable to encode event metadata to JSON.", http.StatusInternalServerError, err)
		return
	}

	alreadyExist, err := ec.EventsStore.Exists(evm.SourceID, evm.EventID)

	if err != nil {
		LogHttpError(w, r, "Unable to persist event.", http.StatusInternalServerError, err)
		return
	}

	if alreadyExist {
		w.WriteHeader(http.StatusOK)
		return
	}

	id, err := ec.EventsStore.Save(evm.SourceID, evm.EventID, evm.Type, eventDataJSON, eventMetadataJSON)

	if err != nil {
		LogHttpError(w, r, "Unable to persist event.", http.StatusInternalServerError, err)
		return
	}

	subscribers := ec.HandlersManager.GetChannels(evm.Type)

	if subscribers != nil {
		savedEvent := savedEventResponse{event: evm}
		savedEvent.Sequence = id

		event, err := json.Marshal(savedEvent)
		if err != nil {
			LogHttpError(w, r, "Unable to encode event to JSON for subscribers.", http.StatusInternalServerError, err)
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
