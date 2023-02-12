package middleware

import (
	"github.com/AliSahib998/ms-parking/model"
	"github.com/AliSahib998/ms-parking/util"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var headers = []string{
	"User-Agent",
	"X-Forwarded-For",
	"Request_ID",
}

func RequestParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		requestID := r.Header.Get(model.HeaderKeyRequestID)
		operation := r.RequestURI
		userAgent := r.Header.Get(model.HeaderKeyUserAgent)
		userIP := r.Header.Get(model.HeaderKeyUserIP)

		if len(requestID) == 0 {
			requestID = uuid.New().String()
		}
		fields := log.Fields{}
		addLoggerParam(fields, model.LoggerKeyRequestID, requestID)
		addLoggerParam(fields, model.LoggerKeyOperation, operation)
		addLoggerParam(fields, model.LoggerKeyUserAgent, userAgent)
		addLoggerParam(fields, model.LoggerKeyUserIP, userIP)

		logger := log.WithFields(fields)
		header := http.Header{}

		for _, v := range headers {
			header.Add(v, r.Header.Get(v))
		}

		ctx = util.WithLogger(ctx, logger)
		ctx = util.WithHeader(ctx, header)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addLoggerParam(fields log.Fields, field string, value string) {
	if len(value) > 0 {
		fields[field] = value
	}
}
