package server

import (
	"fmt"
	"net/http"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/service"
)

func NewHTTPServer(media *service.MediaService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "OK")
	})
	media.RegisterRoutes(mux)
	return mux
}
