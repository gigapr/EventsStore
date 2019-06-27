package main

import (
	"fmt"
	"time"
)

// HTTP/1.1 200 OK
// Content-Type: application/hal+json
// Content-Length: 784

// {
//   "_embedded" : {
//     "jobSearchResultList" : [ {
//       "id" : "425c4a6a-a069-4849-a66e-d08d6e8d1912",
//       "name" : "List * ... Directories bash job",
//       "user" : "genie",
//       "status" : "SUCCEEDED",
//       "started" : "2017-01-12T18:43:42.566Z",
//       "finished" : "2017-01-12T18:43:42.597Z",
//       "clusterName" : "Local laptop",
//       "commandName" : "Unix Bash command",
//       "runtime" : "PT0.031S",
//       "_links" : {
//         "self" : {
//           "href" : "https://genie.example.com/api/v3/jobs/425c4a6a-a069-4849-a66e-d08d6e8d1912"
//         }
//       }
//     } ]
//   },
//   "_links" : {
//     "self" : {
//       "href" : "https://genie.example.com/api/v3/jobs?user=genie"
//     }
//   },
//   "page" : {
//     "size" : 10,
//     "totalElements" : 1,
//     "totalPages" : 1,
//     "number" : 0
//   }
// }
// ```

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
	Number        int `json:"number"`
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

func newEventsResponse(events []savedEventResponse, links map[string]Link, size int, totalElements int, totalPages int, number int) eventsResponse {
	return eventsResponse{
		Embedded: embedded{
			EventsList: events,
		},
		Links: links,
		Page: page{
			Size:          size,
			TotalElements: totalElements,
			TotalPages:    totalPages,
			Number:        number,
		},
	}
}

func newLinks(startFrom int, eventType string, totalNumberOfRecords int) map[string]Link {
	links := map[string]Link{}
	template := "/events/%d/%s"
	if startFrom > 1 {
		links["first"] = Link{
			Href: fmt.Sprintf(template, 1, eventType), //"The first page for this search"
		}
	}

	if startFrom > 1 {
		n := 1
		if (startFrom - pageSize) > 1 {
			n = startFrom - pageSize
		}
		links["prev"] = Link{
			Href: fmt.Sprintf(template, n, eventType), //"The previous page for this search"
		}
	}

	if (startFrom + pageSize) <= totalNumberOfRecords {
		links["next"] = Link{
			Href: fmt.Sprintf(template, startFrom+pageSize, eventType), //"The next page for this search"
		}
	}

	if (totalNumberOfRecords - pageSize) > startFrom {
		links["last"] = Link{
			Href: fmt.Sprintf(template, totalNumberOfRecords-pageSize, eventType), //"The last page for this search"
		}
	}

	return links
}
