[< Back](../README.md)

# Storage

To run the Postgres database execute from the current directory

```

docker build -t eventsDb .

docker run -p 5432:5432 eventsDb

```

Sql scripts placed in the sql directory are automatically applied to the Postgres instance in the docker container on build