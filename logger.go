package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	appLogger  = logging.NewLogger()
	httpLogger = newAccessLogger()
)

func newAccessLogger() func(next http.Handler) http.Handler {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	return middleware.RequestLogger(&accessLogger{logger})
}

type accessLogger struct {
	Logger *logrus.Logger
}

func (l *accessLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	ctx := r.Context()
	entry := &accessLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		logFields["req_id"] = reqID
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["remote_addr"] = r.RemoteAddr
	logFields["xff"] = r.Header.Get("X-Forwarded-For")
	logFields["user_agent"] = r.UserAgent()
	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	logFields["path"] = r.URL.Path

	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		logFields[logging.FieldKeyTraceID] = span.Context().TraceID()
		logFields[logging.FieldKeySpanID] = span.Context().SpanID()
	}

	entry.Logger = entry.Logger.WithFields(logFields)
	return entry
}

type accessLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *accessLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed_ms":   float64(elapsed.Nanoseconds()) / 1000000.0,
	})
	l.Logger.Infoln("request complete")
}

func (l *accessLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
