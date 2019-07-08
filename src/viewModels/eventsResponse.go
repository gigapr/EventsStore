package viewModels

import (
	"fmt"
	"time"
)

type EventsResponse struct {
	Embedded Embedded        `json:"_embedded"`
	Links    map[string]Link `json:"_links"` //this is not serialising properly
	Page     Page            `json:"page"`
}

type Page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	// Number        int `json:"number"`
}

type Link struct {
	Href string `json:"href"`
}

type Embedded struct {
	EventsList []SavedEventResponse `json:"eventsList"`
}

type SavedEventResponse struct {
	Sequence int       `json:"sequence"`
	Received time.Time `json:"received"`
	Event
}

type Event struct {
	SourceID string      `json:"sourceId" validate:"required"`
	EventID  string      `json:"eventId" validate:"required"`
	Type     string      `json:"type" validate:"required"`
	Data     interface{} `json:"data" validate:"required"`
	Metadata interface{} `json:"metadata"`
}

func NewLinks(startFrom int, eventType string, sourceID string, first int, last int, totalNumberOfRecords int, pageSize int) map[string]Link {
	links := map[string]Link{}
	template := "%s/events/%s/%d"

	// links["first"] = Link{
	// 	Href: fmt.Sprintf(template, sourceID, eventType, first), //"The first page for this search"
	// }

	if startFrom > 1 {
		n := 1
		if (startFrom - pageSize) > 1 {
			n = startFrom - pageSize
		}
		links["prev"] = Link{
			Href: fmt.Sprintf(template, sourceID, eventType, n), //"The previous page for this search"
		}
	}

	if (startFrom + pageSize) <= totalNumberOfRecords {
		links["next"] = Link{
			Href: fmt.Sprintf(template, sourceID, eventType, startFrom+pageSize), //"The next page for this search"
		}
	}

	// if (totalNumberOfRecords - pageSize) > startFrom {
	// 	links["last"] = Link{
	// 		Href: fmt.Sprintf(template, sourceID, eventType, last), //"The last page for this search"
	// 	}
	// }

	return links
}
