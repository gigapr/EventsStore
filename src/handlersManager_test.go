package main

import "testing"

func Test_can_subscribe_to_a_topic(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"
	channel := handlersManager.Subscribe(topic)

	if channel == nil {
		t.Error("HandlersManager should create channel")
	}

	handlers := len(handlersManager.handlers)
	if handlers != 1 {
		t.Errorf("HandlersManager should have 1 handlers, got: %d.", handlers)
	}

	channels := len(handlersManager.handlers[topic])
	if channels != 1 {
		t.Errorf("HandlersManager should have 1 channel for '%s', got: %d.", topic, channels)
	}
}

func Test_when_subscribing_to_a_topic_should_handle_multiple_channels_on_the_same_topic(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"
	handlersManager.Subscribe(topic)
	handlersManager.Subscribe(topic)

	handlers := len(handlersManager.handlers)
	if handlers != 1 {
		t.Errorf("HandlersManager should have 1 channel, got: %d.", handlers)
	}

	channels := len(handlersManager.handlers[topic])
	if channels != 2 {
		t.Errorf("HandlersManager should have 2 channels for '%s', got: %d.", topic, channels)
	}
}

func Test_can_unsubscribe_channel_from_a_topic(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"
	channel := handlersManager.Subscribe(topic)

	handlersManager.Unsubscribe(topic, channel)

	handlers := len(handlersManager.handlers)
	if handlers != 0 {
		t.Errorf("HandlersManager should have 0 channel, got: %d.", handlers)
	}
}

func Test_when_unsubscribing_a_channel_doesnt_remove_handlers_if_there_are_subscribed_channels(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"
	handlersManager.Subscribe(topic)
	channel := handlersManager.Subscribe(topic)

	handlersManager.Unsubscribe(topic, channel)

	handlers := len(handlersManager.handlers)
	if handlers != 1 {
		t.Errorf("HandlersManager should have 1 channel, got: %d.", handlers)
	}

	channels := handlersManager.handlers[topic]
	channelsLengths := len(channels)

	if channelsLengths != 1 {
		t.Errorf("HandlersManager should have 1 channels for '%s', got: %d.", topic, channelsLengths)
	}

	if channels[0] == channel {
		t.Error("Unsubscibed incorrect channel")
	}
}

func Test_can_get_channels_for_a_topic(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"
	handlersManager.Subscribe(topic)
	handlersManager.Subscribe(topic)
	handlersManager.Subscribe(topic)
	handlersManager.Subscribe(topic)

	channelsLengths := len(handlersManager.GetChannels(topic))

	if channelsLengths != 4 {
		t.Errorf("HandlersManager should have 4 channels for '%s', got: %d.", topic, channelsLengths)
	}
}

func Test_retunrs_nil_if_there_areent_any_channels_for_a_topic(t *testing.T) {
	handlersManager := NewHandlersManager()

	topic := "SomeTopic"

	channelsLengths := handlersManager.GetChannels(topic)

	if channelsLengths != nil {
		t.Errorf("HandlersManager should not have any channels for '%s', got: %d.", topic, channelsLengths)
	}
}
