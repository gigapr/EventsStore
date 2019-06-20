package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func LogHttpError(w http.ResponseWriter, r *http.Request, errorMessage string, httpStatusCode int, err error) {
	log := log.WithFields(logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"message":         errorMessage,
	})

	log.Error(r, errorMessage, err)

	http.Error(w, errorMessage, httpStatusCode)
}
