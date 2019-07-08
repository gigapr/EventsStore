package main

import (
	"errors"

	"github.com/sirupsen/logrus"
)

func logFatal(message string, err ...error) {
	if err == nil {
		err = []error{errors.New(message)}
	}
	log.WithFields(logrus.Fields{"message": message}).Fatal(err)
}
