package main

func mapEvents(eventsDto []eventDto) []savedEventResponse {
	events := []savedEventResponse{}
	for i := range eventsDto {
		current := eventsDto[i]
		events = append(events, savedEventResponse{
			Sequence: current.ID,
			Received: current.Received,
			event: event{
				SourceID: current.SourceID,
				Type:     current.EventType,
				EventID:  current.EventID,
				Data:     current.EventData,
				Metadata: current.Metadata,
			},
		})
	}
	return events
}
