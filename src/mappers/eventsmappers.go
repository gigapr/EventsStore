package mappers

import (
	"gigapr/eventsstore/persistence"
	"gigapr/eventsstore/viewModels"
)

func MapEvents(eventsDto []persistence.EventDto) []viewModels.SavedEventResponse {
	events := []viewModels.SavedEventResponse{}
	for i := range eventsDto {
		current := eventsDto[i]
		events = append(events, viewModels.SavedEventResponse{
			Sequence: current.ID,
			Received: current.Received,
			Event: viewModels.Event{
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
