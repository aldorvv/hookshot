package webhook

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	database "echook.io/pkg/db"
)

type Handler struct {
	db *database.Database
}

func NewHandler(db *database.Database) *Handler {
	return &Handler{db: db}
}

func JSON[T any](r *http.Request) (T, error) {
	var zero T
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		return zero, errors.New("could not deserialize")
	}
	defer r.Body.Close()

	var v T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&v); err != nil {
		return zero, err
	}

	return v, nil
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	input, err := JSON[*WebhookInput](r)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if err := Create(h.db, input); err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
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
