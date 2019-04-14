Events Store

Prerequisites
- Mongo DB `docker run -d -p 27017:27017 mongo`

To run the application execute

```
make all
make run
```

To run the application within a Docker container 

```
docker build -t eventstore .
docker run -p 4000:4000 eventstore 
```

Events can be Posted to `/event`

```
{
	"type": "userCreated",
	"data": "{ name: 'Gaetano', surname: 'Santonastaso' }",
	"sourceId": "sourceId"
}
```

To receive new events notifications subscribe via WebSocket to /subscribe
