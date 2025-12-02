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
		w.Write([]byte(`{"message": "hearthbeat!"}`))
	})

	db, err := database.NewDatabase("database.sql")
	if err != nil {
		panic(err)
	}

	webhookHandler := webhook.NewHandler(db)
	r.Post("/webhooks", webhookHandler.Post)
	r.Get("/webhooks", webhookHandler.List)

	srv := &http.Server{
		Addr:        fmt.Sprintf("0.0.0.0:%d", 2407),
		Handler:     r,
		IdleTimeout: 60 * time.Second,
	}
	srv.ListenAndServe()
}
