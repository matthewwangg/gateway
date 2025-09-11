package middleware

import (
	"log"
	"net/http"
	"strconv"

	logger "github.com/matthewwangg/gateway/internal/logger"
	metrics "github.com/matthewwangg/gateway/internal/metrics"
)

func init() {
	log.SetFlags(0)
}

type LogResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (lrw *LogResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.w.WriteHeader(code)
}

func (lrw *LogResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *LogResponseWriter) Write(b []byte) (int, error) {
	return lrw.w.Write(b)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &LogResponseWriter{w: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		metrics.Tracker.RecordRequest(r.URL.Path, lrw.statusCode)

		logger.Log.Info(r.Method + " " + r.URL.Path + " " + strconv.Itoa(lrw.statusCode))
	})
}
