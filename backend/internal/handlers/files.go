package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saint0x/file-storage-app/backend/internal/services/storage"
	"github.com/saint0x/file-storage-app/backend/internal/services/websocket"
)

func UploadFile(b2Service *storage.B2Service, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate a unique key for the file
		key := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)

		err = b2Service.UploadFile(r.Context(), key, file)
		if err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			return
		}

		// After successful upload
		fileInfo := map[string]string{"key": key, "name": header.Filename}
		hub.BroadcastUpdate("file_uploaded", fileInfo)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"key": key})
	}
}
