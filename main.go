package main

import (
	"fmt"
	"net/http"
	"time"

	database "echook.io/pkg/db"
	"echook.io/pkg/webhook"
	"github.com/go-chi/chi/v5"
)

type App struct {
	db *database.Database
}

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "heartbeat!"}`))
	})

	// TODO: Replace by env var
	db, err := database.NewDatabase("echook.db")
	if err != nil {
		panic(err)
	}

	webhookHandler := webhook.NewHandler(db)
	r.HandleFunc("/w/{endpoint}", webhookHandler.Capture)
	r.Get("/api/webhooks", webhookHandler.List)

	srv := &http.Server{
		// IDEM, replace port by an env var
		Addr:        fmt.Sprintf("0.0.0.0:%d", 2407),
		Handler:     r,
		IdleTimeout: 60 * time.Second,
	}
	srv.ListenAndServe()
}
