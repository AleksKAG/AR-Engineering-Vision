package store

import (
	"context"
	"time"
	"github.com/google/uuid"
)

type Issue struct {
	ID             string    `json:"id"`
	ElementID      string    `json:"element_id"`
	PhotoURL       string    `json:"photo_url"`
	Comment        string    `json:"comment"`
	DeviationMM    float64   `json:"deviation_mm"`
	AIDetectedType string    `json:"ai_detected_type"`
	AIConfidence   float64   `json:"ai_confidence"`
	IsMatch        bool      `json:"is_match"`
	Status         string    `json:"status"` // "open", "fixed"
	CreatedAt      time.Time `json:"created_at"`
}

func (d *DB) CreateIssueWithAI(ctx context.Context, elementID, photoURL, comment string, deviation float64, aiType string, aiConf float64, isMatch bool) (*Issue, error) {
	id := uuid.New().String()
	_, err := d.pool.Exec(ctx, `
		INSERT INTO issues (id, element_id, photo_url, comment, deviation_mm, ai_detected_type, ai_confidence, is_match, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'open')`,
		id, elementID, photoURL, comment, deviation, aiType, aiConf, isMatch)
	if err != nil {
		return nil, err
	}
	return &Issue{
		ID: id, ElementID: elementID, PhotoURL: photoURL, Comment: comment,
		DeviationMM: deviation, AIDetectedType: aiType, AIConfidence: aiConf,
		IsMatch: isMatch, Status: "open", CreatedAt: time.Now(),
	}, nil
}

// GetIssuesByProject нужен для генерации PDF
func (d *DB) GetIssuesByProject(ctx context.Context, projectID string) ([]Issue, error) {
	query := `
		SELECT i.id, i.element_id, i.photo_url, i.comment, i.deviation_mm, i.ai_detected_type, i.ai_confidence, i.is_match, i.status, i.created_at 
		FROM issues i
		JOIN elements e ON i.element_id = e.id
		JOIN rooms r ON e.room_id = r.id
		WHERE r.project_id = $1
		ORDER BY i.created_at DESC`
	
	rows, err := d.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []Issue
	for rows.Next() {
		var i Issue
		err := rows.Scan(&i.ID, &i.ElementID, &i.PhotoURL, &i.Comment, &i.DeviationMM, &i.AIDetectedType, &i.AIConfidence, &i.IsMatch, &i.Status, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		issues = append(issues, i)
	}
	return issues, nil
}
