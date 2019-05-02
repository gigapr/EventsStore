package main

import "log"

type HandlersManager struct {
	handlers map[string][]chan []byte
}

func NewHandlersManager() *HandlersManager {

	es := new(HandlersManager)
	es.handlers = make(map[string][]chan []byte)
	return es
}

func (hm HandlersManager) Unregister(topic string, channel chan []byte) {

	for i, other := range hm.handlers[topic] {
		if other == channel {
			hm.handlers[topic] = append(hm.handlers[topic][:i], hm.handlers[topic][i+1:]...)
			log.Println("Channel unregistered", channel)
			break
		}
	}
}

func (hm HandlersManager) Register(topic string) chan []byte {
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

func (hm HandlersManager) Get(topic string) []chan []byte {
	if _, ok := hm.handlers[topic]; ok {
		return hm.handlers[topic]
	}
	return nil
}
