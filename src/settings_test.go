package main

import (
	"os"
	"testing"
)

func Test_can_initialise_settings_with_default_value(t *testing.T) {
	settings := InitialiseSettings()

	if settings.Port != "4000" {
		t.Errorf("Port should have 4000 as default vaule, got: %s.", settings.Port)
	}
}

func Test_can_read_overridden_port_value(t *testing.T) {
	environmentVariable := "1000"
	os.Setenv("port", environmentVariable)

	settings := InitialiseSettings()

	if settings.Port != environmentVariable {
		t.Errorf("Port should have %s as default vaule, got: %s.", environmentVariable, settings.Port)
	}
}
