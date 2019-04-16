package main

type HandlersManager struct {
	handlers map[string]chan []byte
}

func NewHandlersManager() *HandlersManager {

	es := new(HandlersManager)
	es.handlers = make(map[string]chan []byte)
	return es
}

func (hm HandlersManager) Get(topic string) chan []byte {
	if _, ok := hm.handlers[topic]; ok {
	} else {
		hm.handlers[topic] = make(chan []byte)
	}

	return hm.handlers[topic]
}
