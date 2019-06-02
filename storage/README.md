[< Back](../README.md)

# Storage

To run the Postgres database execute from the current directory

```

docker build -t eventsdb .

docker run -p 5432:5432 eventsdb

```

Sql scripts placed in the sql directory are automatically applied to the Postgres instance in the docker container on build