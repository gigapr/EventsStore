package websocket //import "gigapr/eventsstore/websocket"

//HandlersManager is responsible for managing clients connections
type HandlersManager struct {
	Handlers map[string][]chan []byte
}

//Subscribe creates a channel for a particular topic
func (hm HandlersManager) Subscribe(topic string) chan []byte {
	channel := make(chan []byte)
	if _, ok := hm.Handlers[topic]; ok {
		hm.Handlers[topic] = append(hm.Handlers[topic], channel)
	} else {
		hm.Handlers[topic] = []chan []byte{
			channel,
		}
	}

	return channel
}

//Unsubscribe is called when client doesn't want to receive events or connection is broken
func (hm HandlersManager) Unsubscribe(topic string, channel chan []byte) {

	for i, other := range hm.Handlers[topic] {
		if other == channel {
			hm.Handlers[topic] = append(hm.Handlers[topic][:i], hm.Handlers[topic][i+1:]...)
			// log.Debug("Channel unregistered", channel)
			if len(hm.Handlers[topic]) == 0 {
				delete(hm.Handlers, topic)
				// log.Debug(fmt.Sprintf("Removed handler for topic '%s'", topic))
			}
			break
		}
	}
}

//GetChannels returns all the channels subscribed to a particular topic
func (hm HandlersManager) GetChannels(topic string) []chan []byte {
	if _, ok := hm.Handlers[topic]; ok {
		return hm.Handlers[topic]
	}
	return nil
}

//GetAllChannels returns all connected channels
func (hm HandlersManager) GetAllChannels() map[string][]chan []byte {
	return hm.Handlers
}
