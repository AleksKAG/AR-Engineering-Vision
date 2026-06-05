package store

import (
	"context"
	"encoding/json"
)

// Room соответствует таблице rooms
type Room struct {
	ID        string          `json:"id"`
	ProjectID string          `json:"project_id"`
	Name      string          `json:"name"`
	Bbox      json.RawMessage `json:"bbox"`
}

// ElementInput используется для вставки новых элементов
type ElementInput struct {
	ID          string          `json:"id"`
	RoomID      string          `json:"room_id"`
	Type        string          `json:"type"`
	WorldCoords json.RawMessage `json:"world_coords"`
	Properties  json.RawMessage `json:"properties"`
}

// SaveBimData атомарно сохраняет комнаты и элементы
func (d *DB) SaveBimData(ctx context.Context, projectID, glbS3Key string, rooms []Room, elements []ElementInput) error {
	// Начинаем транзакцию
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Обновляем запись в models, указывая, что glb готов
	_, err = tx.Exec(ctx, `UPDATE models SET s3_key = $1 WHERE project_id = $2`, glbS3Key, projectID)
	if err != nil {
		return err
	}

	// 2. Вставляем комнаты
	for _, room := range rooms {
		_, err = tx.Exec(ctx, 
			`INSERT INTO rooms (id, project_id, name, bbox) VALUES ($1, $2, $3, $4)`,
			room.ID, room.ProjectID, room.Name, room.Bbox)
		if err != nil {
			return err
		}
	}

	// 3. Вставляем элементы
	for _, elem := range elements {
		_, err = tx.Exec(ctx, 
			`INSERT INTO elements (id, room_id, type, world_coords, properties) VALUES ($1, $2, $3, $4, $5)`,
			elem.ID, elem.RoomID, elem.Type, elem.WorldCoords, elem.Properties)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
