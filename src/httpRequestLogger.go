package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPRequestLogger struct {
	*logrus.Logger
}

func (rl *HTTPRequestLogger) Error(r *http.Request, message string, err error) {
	rl.decorate(r, message).Error(err)
}

func (rl *HTTPRequestLogger) Debug(r *http.Request, message string, err error) {
	rl.decorate(r, message).Debug(err)
}

func (rl *HTTPRequestLogger) decorate(r *http.Request, message string) *logrus.Entry {
	log := rl.WithFields(logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"message":         message,
	})

	return log
}
