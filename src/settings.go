package main

import (
	"log"
	"os"
	"strconv"
)

//Settings type that holds application setttings
type Settings struct {
	Port             string
	DatabaseHost     string
	DatabasePort     int
	DatabaseUsername string
	DatabasePassword string
	DatabaseName     string
}

/*InitialiseSettings initialise a new Settings.
  Configurations can be overridden via environment variables.
  Available settings:
  [ port, databaseHost, databsePort, databaseUsername, databasePassword, databaseName ]
*/
func InitialiseSettings() *Settings {
	settings := new(Settings)
	settings.Port = getEnv("port", "4000")
	settings.DatabaseHost = getEnv("databaseHost", "localhost")

	databasePort, err := strconv.Atoi(getEnv("databsePort", "5432"))
	if err != nil {
		log.Fatal(err)
	}

	settings.DatabasePort = databasePort
	settings.DatabaseUsername = getEnv("databaseUsername", "postgressuperuser")
	settings.DatabasePassword = getEnv("databasePassword", "mysecretpassword")
	settings.DatabaseName = getEnv("databaseName", "eventsStore")

	return settings
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
