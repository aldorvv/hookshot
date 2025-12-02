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
	body, _ := io.ReadAll(r.Body)
	headers, _ := json.Marshal(r.Header)

	input := &WebhookInput{
		Endpoint: endpoint,
		Method:   r.Method,
		Headers:  headers,
		Body:     body,
		IP:       r.RemoteAddr,
	}

	if err := Create(h.db, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "captured"}`))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := Get(h.db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(record)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	results, err := List(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
