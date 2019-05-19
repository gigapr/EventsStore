package main

import (
	"os"
)

//Settings type that holds application setttings
type Settings struct {
	Port string
}

/*InitialiseSettings initialise a new Settings.
  Configurations can be overridden via environment variables.
  Available settings:
  [ port ]
*/
func InitialiseSettings() *Settings {
	settings := new(Settings)
	settings.Port = getEnv("port", "4000")
	return settings
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
