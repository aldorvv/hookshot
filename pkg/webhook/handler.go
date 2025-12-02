package webhook

import (
	"encoding/json"
	"io"
	"net/http"

	database "echook.io/pkg/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Capture(w http.ResponseWriter, r *http.Request) {
	endpoint := chi.URLParam(r, "endpoint")
	bodyBytes, _ := io.ReadAll(r.Body)
	headersJSON, _ := json.Marshal(r.Header)

	input := &WebhookInput{
		Endpoint: endpoint,
		Method:   r.Method,
		Headers:  headersJSON,
		Body:     bodyBytes,
		IP:       r.RemoteAddr,
	}

	if err := Create(h.db, input); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(`{"status": "captured"}`))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	results, err := List(h.db)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}
