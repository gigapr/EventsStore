[![Build Status](https://travis-ci.com/gigapr/EventsStore.svg?branch=master)](https://travis-ci.com/gigapr/
EventsStore)

# Events Store

To run the application execute

```
make all
make run
```

or  
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
  "type": "userCreated",
  "data": "{ name: 'Gaetano', surname: 'Santonastaso' }",
  "sourceId": "sourceId"
}
```

To receive new events notifications subscribe via WebSocket to `/subscribe`

There is an example client in the client directory

