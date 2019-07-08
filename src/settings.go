package main

import (
	"fmt"
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
	PageSize         int
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

	databasePortString := getEnv("databsePort", "5432")
	databasePort, err := strconv.Atoi(databasePortString)
	if err != nil {
		logFatal(fmt.Sprintf("Unable to convert '%s' to int", databasePortString), err)
	}

	settings.DatabasePort = databasePort
	settings.DatabaseUsername = getEnv("databaseUsername", "postgressuperuser")
	settings.DatabasePassword = getEnv("databasePassword", "mysecretpassword")
	settings.DatabaseName = getEnv("databaseName", "eventsStore")
	settings.PageSize = 2

	return settings
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
