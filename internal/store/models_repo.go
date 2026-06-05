package store

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type Model struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	S3Key       string    `json:"s3_key"`
	CreatedAt   time.Time `json:"created_at"`
}

func (d *DB) CreateModel(ctx context.Context, projectID, filename, contentType, s3Key string) (*Model, error) {
	id := uuid.New().String()
	_, err := d.pool.Exec(ctx, 
		`INSERT INTO models (id, project_id, filename, content_type, s3_key) VALUES ($1,$2,$3,$4,$5)`, 
		id, projectID, filename, contentType, s3Key)
	if err != nil {
		return nil, err
	}
	return &Model{ID: id, ProjectID: projectID, Filename: filename, ContentType: contentType, S3Key: s3Key, CreatedAt: time.Now()}, nil
}
