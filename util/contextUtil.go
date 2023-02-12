package util

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	contextLogger = "contextLogger"
	contextHeader = "contextHeader"
)

// WithLogger add logger to context
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, contextLogger, logger)
}

// WithHeader add header to context
func WithHeader(ctx context.Context, header http.Header) context.Context {
	return context.WithValue(ctx, contextHeader, header)
}

// GetLogger get logger from context
func GetLogger(ctx context.Context) *logrus.Entry {
	logger, ok := ctx.Value(contextLogger).(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(logrus.New())
	}
	return logger
}

// GetHeader get header from context
func GetHeader(ctx context.Context) http.Header {
	header, ok := ctx.Value(contextHeader).(http.Header)
	if !ok {
		return http.Header{}
	}
	return header
}
