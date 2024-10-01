package apimiddleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/saint0x/file-storage-app/backend/pkg/logger"
)

func RequestLogger(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			latency := time.Since(start)

			log.Info("Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"latency", latency,
				"bytes", ww.BytesWritten(),
			)
		})
	}
}

func Recoverer(log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.Error("Panic recovered",
						"error", rvr,
						"stack", string(debug.Stack()),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// WrapResponseWriter is a wrapper around http.ResponseWriter that allows us to capture the status code and bytes written
type WrapResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int) *WrapResponseWriter {
	return &WrapResponseWriter{ResponseWriter: w}
}

func (w *WrapResponseWriter) Status() int {
	return w.status
}

func (w *WrapResponseWriter) BytesWritten() int {
	return w.bytes
}

func (w *WrapResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *WrapResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
