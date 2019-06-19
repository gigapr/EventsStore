package main

import "time"

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
	SourceID string      `json:"sourceId"`
	EventID  string      `json:"eventId"`
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
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

func newLinks() map[string]Link {
	return map[string]Link{
		"self": Link{
			Href: "https://genie.example.com/api/v3/jobs/425c4a6a-a069-4849-a66e-d08d6e8d1912",
		},
		"first": Link{
			Href: "The first page for this search"},
		"prev": Link{
			Href: "The previous page for this search"},
		"next": Link{
			Href: "The next page for this search"},
		"last": Link{
			Href: "The last page for this search"},
	}
}
