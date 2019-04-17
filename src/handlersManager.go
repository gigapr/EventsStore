package main

type HandlersManager struct {
	handlers map[string][]chan []byte
}

func NewHandlersManager() *HandlersManager {

	es := new(HandlersManager)
	es.handlers = make(map[string][]chan []byte)
	return es
}

// func (hm HandlersManager) Unregister(topic string) []chan []byte {

// }

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
