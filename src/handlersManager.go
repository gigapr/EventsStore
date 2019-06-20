package main

import (
	"fmt"
)

//HandlersManager is responsible for managing clients connections
type HandlersManager struct {
	handlers map[string][]chan []byte
}

//NewHandlersManager inititalise a new HandlersManager
func NewHandlersManager() *HandlersManager {
	es := new(HandlersManager)
	es.handlers = make(map[string][]chan []byte)
	return es
}

//Subscribe creates a channel for a particular topic
func (hm HandlersManager) Subscribe(topic string) chan []byte {
	channel := make(chan []byte)
	if _, ok := hm.handlers[topic]; ok {
		hm.handlers[topic] = append(hm.handlers[topic], channel)
	} else {
		hm.handlers[topic] = []chan []byte{
			channel,
		}
	}

	return channel
}

//Unsubscribe is called when client doesn't want to receive events or connection is broken
func (hm HandlersManager) Unsubscribe(topic string, channel chan []byte) {

	for i, other := range hm.handlers[topic] {
		if other == channel {
			hm.handlers[topic] = append(hm.handlers[topic][:i], hm.handlers[topic][i+1:]...)
			log.Debug("Channel unregistered", channel)
			if len(hm.handlers[topic]) == 0 {
				delete(hm.handlers, topic)
				log.Debug(fmt.Sprintf("Removed handler for topic '%s'", topic))
			}
			break
		}
	}
}

//GetChannels returns all the channels subscribed to a particular topic
func (hm HandlersManager) GetChannels(topic string) []chan []byte {
	if _, ok := hm.handlers[topic]; ok {
		return hm.handlers[topic]
	}
	return nil
}

//GetAllChannels returns all connected channels
func (hm HandlersManager) GetAllChannels() map[string][]chan []byte {
	return hm.handlers
}
