package webhook

import (
	"encoding/json"
	"time"

	database "echook.io/pkg/db"
	"github.com/google/uuid"
)

type WebhookInput struct {
	Endpoint string          `json:"endpoint"`
	Method   string          `json:"method"`
	Headers  json.RawMessage `json:"headers"`
	Body     json.RawMessage `json:"body"`
	IP       string          `json:"ip"`
}

type WebhookRecord struct {
	ID        string          `json:"id"`
	Endpoint  string          `json:"endpoint"`
	Method    string          `json:"method"`
	Headers   json.RawMessage `json:"headers"`
	Body      json.RawMessage `json:"body"`
	IP        string          `json:"ip"`
	CreatedAt time.Time       `json:"created_at"`
}

func Create(db *database.Database, input *WebhookInput) error {
	id := uuid.New()
	return db.Exec(
		`INSERT INTO webhooks (id, endpoint, method, headers, body, ip, created_at) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		id,
		input.Endpoint,
		input.Method,
		input.Headers,
		input.Body,
		input.IP,
		time.Now(),
	)
}

func Get(db *database.Database, id string) (*WebhookRecord, error) {
	var output WebhookRecord
	row := db.QueryRow(`SELECT id, endpoint, method, headers, body, ip, created_at FROM webhooks WHERE id = ?;`, id)

	return &output, row.Scan(&output.ID, &output.Endpoint, &output.Method, &output.Headers, &output.Body, &output.IP, &output.CreatedAt)
}

func List(db *database.Database) ([]*WebhookRecord, error) {
	results := []*WebhookRecord{}
	rows, err := db.Query(`SELECT id, endpoint, method, headers, body, ip, created_at FROM webhooks;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var record WebhookRecord
		err := rows.Scan(
			&record.ID,
			&record.Endpoint,
			&record.Method,
			&record.Headers,
			&record.Body,
			&record.IP,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, &record)
	}
	return results, nil
}
