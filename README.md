# Events Store | [![Build Status](https://travis-ci.com/gigapr/EventsStore.svg?branch=master)](https://travis-ci.com/gigapr/EventsStore) | ![Docker Pulls](https://img.shields.io/docker/pulls/threeamigos/eventstore.svg)

An experimental Events Store written in Golang that allows to subscribe to events of a particular topic using [WebSocket](https://en.wikipedia.org/wiki/WebSocket). 

## Dependencies

- [Events Database](./persistence/README.md)

To build the application execute from the root directory run:

```
make all
```

To execute the tests run:

```
make test
```


To run the application from the src directory run:
```
go run $(ls -1 *.go | grep -v _test.go)
```


To run the application within a Docker container run:

```
docker build -t eventstore .

docker run -p 4000:4000 eventstore 

```

To POST Events: 

```
curl -X POST \
  http://localhost:4000/event \
  -H 'content-type: application/json' \
  -d '{
        "type": "userCreated",
        "data": "{ name: 'Gaetano', surname: 'Santonastaso' }",
        "sourceId": "sourceId",
        "eventId": "eventId",
        "metadata": "{ key: 'value' }"
      }'
```

To receive new events notifications subscribe via WebSocket to `/subscribe?topic=eventType`

To run the [client example](./client/client.go) from the client directory run: 

```
go run client.go
```

Pre built Docker image can be downloaded from Dockerhub at [eventsstore](https://cloud.docker.com/u/threeamigos/repository/docker/threeamigos/eventstore)
