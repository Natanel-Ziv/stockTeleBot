package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type responseData struct {
    status int
    size int
}

type loggingResponseWriter struct {
    http.ResponseWriter
    responseData *responseData
}

func(r *loggingResponseWriter) Write(b []byte) (int, error) {
    size, err := r.ResponseWriter.Write(b)
    r.responseData.size += size
    return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
    r.ResponseWriter.WriteHeader(statusCode)
    r.responseData.status = statusCode
}

func WithLogging(h http.Handler, logger *logrus.Logger) http.Handler {
    loggingFn := func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        responseData := &responseData{}
        lrw := loggingResponseWriter{
            ResponseWriter: w,
            responseData: responseData,
        }

        h.ServeHTTP(&lrw, r)

        duration := time.Since(start)
        logger.WithFields(logrus.Fields{
            "uri": r.RequestURI,
            "method": r.Method,
            "status": responseData.status,
            "duration": duration.String(),
            "size": responseData.size,
        }).Info("")
    }
    return http.HandlerFunc(loggingFn)
}