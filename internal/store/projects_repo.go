package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (d *DB) CreateProject(ctx context.Context, ownerID, name, description string) (*Project, error) {
	id := uuid.New().String()
	_, err := d.pool.Exec(ctx, `INSERT INTO projects (id, owner_id, name, description) VALUES ($1,$2,$3,$4)`, id, ownerID, name, description)
	if err != nil {
		return nil, err
	}
	return &Project{ID: id, OwnerID: ownerID, Name: name, Description: description, CreatedAt: time.Now()}, nil
}

func (d *DB) ListProjectsByOwner(ctx context.Context, ownerID string) ([]Project, error) {
	rows, err := d.pool.Query(ctx, `SELECT id, owner_id, name, description, created_at FROM projects WHERE owner_id = $1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (d *DB) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	row := d.pool.QueryRow(ctx, `SELECT id, owner_id, name, description, created_at FROM projects WHERE id = $1`, id)
	var p Project
	if err := row.Scan(&p.ID, &p.OwnerID, &p.Name, &p.Description, &p.CreatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}
