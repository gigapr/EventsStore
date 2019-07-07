package main

import (
	"fmt"
	"time"
)

type event struct {
	SourceID string      `json:"sourceId" validate:"required"`
	EventID  string      `json:"eventId" validate:"required"`
	Type     string      `json:"type" validate:"required"`
	Data     interface{} `json:"data" validate:"required"`
	Metadata interface{} `json:"metadata"`
}

type page struct {
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	// Number        int `json:"number"`
}

type embedded struct {
	EventsList []savedEventResponse `json:"eventsList"`
}

type savedEventResponse struct {
	Sequence int       `json:"sequence"`
	Received time.Time `json:"received"`
	event
}

type eventsResponse struct {
	Embedded embedded        `json:"_embedded"`
	Links    map[string]Link `json:"_links"` //this is not serialising properly
	Page     page            `json:"page"`
}

type Link struct {
	Href string `json:"href"`
}

func newLinks(startFrom int, eventType string, sourceID string, first int, last int, totalNumberOfRecords int) map[string]Link {
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
