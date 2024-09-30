package apimiddleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
)

type TimingInfo struct {
	StartTime          time.Time     `json:"start_time"`
	EndTime            time.Time     `json:"end_time"`
	TotalDuration      time.Duration `json:"total_duration"`
	BodyReadDuration   time.Duration `json:"body_read_duration"`
	ProcessingDuration time.Duration `json:"processing_duration"`
	WriteDuration      time.Duration `json:"write_duration"`
}

type LogEntry struct {
	RequestID    string      `json:"request_id"`
	Method       string      `json:"method"`
	Path         string      `json:"path"`
	RemoteAddr   string      `json:"remote_addr"`
	UserAgent    string      `json:"user_agent"`
	Status       int         `json:"status"`
	Latency      float64     `json:"latency"`
	RequestBody  interface{} `json:"request_body,omitempty"`
	ResponseBody interface{} `json:"response_body,omitempty"`
	Timestamp    time.Time   `json:"timestamp"`
	TimingInfo   TimingInfo  `json:"timing_info"`
}

func Logging(wsHub *websocket.Hub) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Read and store the request body
			var requestBody interface{}
			var bodyReadDuration time.Duration
			if r.Body != nil {
				bodyReadStart := time.Now()
				bodyBytes, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				json.Unmarshal(bodyBytes, &requestBody)
				bodyReadDuration = time.Since(bodyReadStart)
			}

			// Create a response writer that captures the response
			buf := &bytes.Buffer{}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ww.Tee(buf)

			// Call the next handler
			processingStart := time.Now()
			next.ServeHTTP(ww, r)
			processingDuration := time.Since(processingStart)

			// Capture response body
			writeStart := time.Now()
			responseBody := buf.Bytes()
			writeDuration := time.Since(writeStart)

			end := time.Now()

			// Create timing info
			timingInfo := TimingInfo{
				StartTime:          start,
				EndTime:            end,
				TotalDuration:      end.Sub(start),
				BodyReadDuration:   bodyReadDuration,
				ProcessingDuration: processingDuration,
				WriteDuration:      writeDuration,
			}

			// Create log entry
			logEntry := LogEntry{
				RequestID:    middleware.GetReqID(r.Context()),
				Method:       r.Method,
				Path:         r.URL.Path,
				RemoteAddr:   r.RemoteAddr,
				UserAgent:    r.UserAgent(),
				Status:       ww.Status(),
				Latency:      timingInfo.TotalDuration.Seconds(),
				RequestBody:  requestBody,
				ResponseBody: json.RawMessage(responseBody),
				Timestamp:    start,
				TimingInfo:   timingInfo,
			}

			// Send log entry to WebSocket clients
			wsMessage, _ := json.Marshal(logEntry)
			wsHub.BroadcastUpdate("log_entry", string(wsMessage))

			// You can also log to a file or database here if needed
		})
	}
}
