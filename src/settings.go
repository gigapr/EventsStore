package main

import (
	"os"
)

type Setting struct {
	Port string
}

func InitialiseSetting() *Setting {
	settings := new(Setting)
	settings.Port = "4000"
	return settings
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
