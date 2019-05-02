[![Build Status](https://travis-ci.com/gigapr/EventsStore.svg?branch=master)](https://travis-ci.com/gigapr/EventsStore)

[Docker hub](https://cloud.docker.com/u/threeamigos/repository/docker/threeamigos/eventstore)

# Events Store

To run the application execute

```
make all
make run 
```

or from the src directory
```
go run $(ls -1 *.go | grep -v _test.go)
```


To run the application within a Docker container 

```
docker build -t eventstore .
docker run -p 4000:4000 eventstore 
```

Events can be Posted to `/subscribe?topic=eventType`

```
{
  "type": "topicName",
  "data": "{ name: 'Gaetano', surname: 'Santonastaso' }",
  "sourceId": "sourceId"
}
```

To receive new events notifications subscribe via WebSocket to `/subscribe?topic=topicName`

There is an example client in the client directory

```
go run client.go

```

