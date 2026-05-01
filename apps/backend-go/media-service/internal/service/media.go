package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arda-labs/arda/apps/backend-go/media-service/internal/biz"
)

type MediaService struct {
	uc *biz.MediaUsecase
}

func NewMediaService(uc *biz.MediaUsecase) *MediaService {
	return &MediaService{uc: uc}
}

type initUploadRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	OwnerID     string `json:"owner_id"`
	Module      string `json:"module"`
}

type initUploadResponse struct {
	Media     *biz.MediaMetadata `json:"media"`
	UploadURL string             `json:"upload_url"`
	ExpiresAt string             `json:"expires_at"`
}

type mediaURLResponse struct {
	Media       *biz.MediaMetadata `json:"media"`
	DownloadURL string             `json:"download_url"`
	ExpiresAt   string             `json:"expires_at"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *MediaService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/media/upload/init", s.initUpload)
	mux.HandleFunc("POST /v1/media/", s.mediaAction)
	mux.HandleFunc("GET /v1/media/", s.mediaAction)
}

func (s *MediaService) initUpload(w http.ResponseWriter, r *http.Request) {
	var req initUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	result, err := s.uc.InitUpload(r.Context(), biz.InitUploadInput{
		Filename:    req.Filename,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
		OwnerID:     req.OwnerID,
		Module:      req.Module,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, initUploadResponse{Media: result.Media, UploadURL: result.UploadURL, ExpiresAt: result.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")})
}

func (s *MediaService) mediaAction(w http.ResponseWriter, r *http.Request) {
	id, suffix, ok := parseMediaPath(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	switch {
	case r.Method == http.MethodPost && suffix == "confirm":
		s.confirmUpload(w, r, id)
	case r.Method == http.MethodGet && suffix == "url":
		s.getMediaURL(w, r, id)
	default:
		writeError(w, http.StatusNotFound, "not found")
	}
}

func (s *MediaService) confirmUpload(w http.ResponseWriter, r *http.Request, id string) {
	media, err := s.uc.ConfirmUpload(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, media)
}

func (s *MediaService) getMediaURL(w http.ResponseWriter, r *http.Request, id string) {
	result, err := s.uc.GetDownloadURL(r.Context(), id, r.URL.Query().Get("owner_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, mediaURLResponse{Media: result.Media, DownloadURL: result.DownloadURL, ExpiresAt: result.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")})
}

func parseMediaPath(path string) (string, string, bool) {
	rest := strings.TrimPrefix(path, "/v1/media/")
	if rest == path || rest == "" {
		return "", "", false
	}
	parts := strings.Split(rest, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
