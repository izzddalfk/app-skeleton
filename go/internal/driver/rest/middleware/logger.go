package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/izzdalfk/app-skeleton/go/internal/shared/logger"
)

// LoggerContext used to embed system logger instance into root of request context.
// So later any service and its dependency can get the logger from the context
// instead of passing the logger into other components as dependency.
func LoggerContext(apiLogger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggerCtx := context.WithValue(r.Context(), logger.CtxKey, apiLogger)
			next.ServeHTTP(w, r.WithContext(loggerCtx))
		})
	}
}

// LogRequest used to log every requests that come into rest API using logger
func LogRequest(apiLogger logger.Logger) func(next http.Handler) http.Handler {
	recoveryHandler := func(reqStats map[string]interface{}) {
		err := recover()
		if err != nil {
			apiLogger.SetFields(
				logger.Field{
					Key:   "stack_trace",
					Value: string(debug.Stack()),
				},
				logger.Field{
					Key:   "err",
					Value: err,
				},
				logger.Field{
					Key:   "http_request",
					Value: reqStats,
				},
			).Error("http.req")

			// TODO: render the error
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			resp := responseLogger(w)
			reqStats := requestLogFields(r)

			// panic recovery handler
			defer recoveryHandler(reqStats)

			next.ServeHTTP(resp, r)

			end := time.Since(start)

			// response statistics of the request
			respStats := map[string]interface{}{
				"header":       headerLogFields(resp.Header()),
				"status":       resp.statusCode,
				"bytes":        len(resp.body),
				"time_elapsed": fmt.Sprintf("%.3fms", float64(end.Nanoseconds())/1000000.0), // in milliseconds
			}

			reqLogField := logger.Field{
				Key:   "http_request",
				Value: reqStats,
			}
			respLogField := logger.Field{
				Key:   "http_response",
				Value: respStats,
			}

			switch c := resp.statusCode; {
			case c < 400:
				apiLogger.Info("http.req", reqLogField, respLogField)
			case c >= 400:
				if resp.body != nil {
					respStats["body"] = string(resp.body)
				}
				if c >= 500 {
					apiLogger.Error("http.req", reqLogField, respLogField)

					return
				}
				apiLogger.Warn("http.req", reqLogField, respLogField)
			}
		})
	}
}

// requestLogFields used to construct structured logs regarding request data
func requestLogFields(r *http.Request) map[string]interface{} {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	requestFields := map[string]interface{}{
		"scheme": scheme,
		"url":    fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI),
		"method": r.Method,
		"path":   r.URL.Path,
		"proto":  r.Proto,
	}
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		requestFields["request_id"] = reqID
	}
	// write logs for request header
	if len(r.Header) > 0 {
		requestFields["header"] = headerLogFields(r.Header)
	}
	if r.Body != nil {
		// read request body
		data, _ := io.ReadAll(r.Body)
		// since we already read the request body, it will return nil on the further reads
		// inside the next handlers, to avoid this we need to reconstruct the request body
		body := io.NopCloser(bytes.NewBuffer(data))
		r.Body = body

		if len(data) > 0 {
			// the request body will be logged within 1024 KB
			if len(data) > 1024000 {
				data = data[:1024000]
			}
			requestFields["body"] = string(data)
		}
	}

	return requestFields
}

// headerLogFields used to construct structured logs regarding request header
func headerLogFields(header http.Header) map[string]string {
	headerFields := map[string]string{}
	for headerKey, headerValue := range header {
		if len(headerValue) == 0 {
			continue
		}
		headerFields[headerKey] = fmt.Sprintf("%s", headerValue)

		// mask values from specific headers
		maskedKey := strings.ToLower(headerKey)
		switch maskedKey {
		case "authorization", "cookie", "set-cookie", "x-api-key":
			headerFields[headerKey] = "[***]"
		}
	}
	return headerFields
}

type responseLog struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func responseLogger(w http.ResponseWriter) *responseLog {
	// default status is 200
	return &responseLog{ResponseWriter: w, statusCode: http.StatusOK, body: nil}
}

// WriteHeader override status and set response value
func (c *responseLog) WriteHeader(status int) {
	c.statusCode = status
	c.ResponseWriter.WriteHeader(status)
}

func (c *responseLog) Write(b []byte) (int, error) {
	c.body = b
	return c.ResponseWriter.Write(b)
}
