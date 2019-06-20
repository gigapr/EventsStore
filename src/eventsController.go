package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type eventsController struct {
	log             *HTTPRequestLogger
	EventsStore     *EventsStore
	HandlersManager *HandlersManager
}

//InitialiseEventsController register all the EventsController routes
func InitialiseEventsController(router *mux.Router, eventStore *EventsStore, handlersManager *HandlersManager) {

	es := new(eventsController)
	es.EventsStore = eventStore
	es.HandlersManager = handlersManager
	es.log = NewHTTPRequestLogger()

	router.HandleFunc("/event", es.saveEventHandler)
	router.HandleFunc("/events/{startFrom}", es.getEventsHandler)
}

//events?
func (ec *eventsController) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	startFrom, err := strconv.Atoi(vars["startFrom"])
	if err != nil {
		http.Error(w, fmt.Sprintf("'%s' is not a valid numnber", vars["startFrom"]), http.StatusBadRequest)
		return
	}

	eventsDto, err := ec.EventsStore.Get(startFrom)
	if err != nil {
		ec.log.Error(r, fmt.Sprintf("Unable to get events starting from %d.", startFrom), err)
		http.Error(w, "Unable to process the request.", http.StatusInternalServerError) ///should have erroras codes
		return
	}

	events := mapEvents(eventsDto)
	links := newLinks()
	eventsList := newEventsResponse(events, links, 1, 1, 1, 1)
	json, err := json.Marshal(eventsList)
	if err != nil {
		ec.log.Error(r, "Unable to encode events response to JSON for subscribers.", err)
		http.Error(w, "Unable to process the request.", http.StatusInternalServerError) ///should have erroras codes
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	eventDataJSON, err := json.Marshal(evm.Data)

	if err != nil {
		message := "Unable to encode event data to JSON."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	eventMetadataJSON, err := json.Marshal(evm.Metadata)
	if err != nil {
		message := "Unable to encode event metadata to JSON."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)

		return
	}

	alreadyExist, err := ec.EventsStore.Exists(evm.SourceID, evm.EventID)

	if err != nil {
		message := "Unable to persist event."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)

		return
	}

	if alreadyExist {
		w.WriteHeader(http.StatusOK)
		return
	}

	id, err := ec.EventsStore.Save(evm.SourceID, evm.EventID, evm.Type, eventDataJSON, eventMetadataJSON)

	if err != nil {
		message := "Unable to persist event."
		ec.log.Error(r, message, err)
		http.Error(w, message, http.StatusInternalServerError)

		return
	}

	subscribers := ec.HandlersManager.GetChannels(evm.Type)

	if subscribers != nil {
		savedEvent := savedEventResponse{event: evm}
		savedEvent.Sequence = id

		event, err := json.Marshal(savedEvent)
		if err != nil {
			message := "Unable to encode event to JSON for subscribers."
			ec.log.Error(r, message, err)
			http.Error(w, message, http.StatusInternalServerError)

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
