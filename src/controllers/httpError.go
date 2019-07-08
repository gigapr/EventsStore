package controllers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func httpError(l *logrus.Logger, w http.ResponseWriter, r *http.Request, errorMessage string, httpStatusCode int, err error) {
	fields := logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"message":         errorMessage,
	}

	log := l.WithFields(fields)

	log.Error(r, errorMessage, err)

	http.Error(w, errorMessage, httpStatusCode)
}
