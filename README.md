[![Build Status](https://travis-ci.com/gigapr/EventsStore.svg?branch=master)](https://travis-ci.com/gigapr/EventsStore)

[Docker hub](https://cloud.docker.com/u/threeamigos/repository/docker/threeamigos/eventstore)

# Events Store

To build the application execute from the root directory

```
make all
```

To run the application from the src directory run
```
go run $(ls -1 *.go | grep -v _test.go)
```


To run the application within a Docker container 

```
docker build -t eventstore .
docker run -p 4000:4000 eventstore 
```

To POST Events 

```
curl -X POST \
  http://localhost:4000/event \
  -H 'acc: application/json' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
        "type": "userCreated",
        "data": "{ name: 'Gaetano', surname: 'Santonastaso' }",
        "sourceId": "sourceId"
      }'
```

To receive new events notifications subscribe via WebSocket to `/subscribe?topic=eventType`

[Client example](./client/client.go)
