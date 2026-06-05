package store

import (
	"context"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Element struct {
	ID          string          `json:"id"`
	RoomID      string          `json:"room_id"`
	Type        string          `json:"type"` // duct, pipe, socket
	WorldCoords json.RawMessage `json:"world_coords"` // JSONB с координатами X,Y,Z
	Properties  json.RawMessage `json:"properties"`
	CreatedAt   time.Time       `json:"created_at"`
}

func (d *DB) GetElementsByRoomID(ctx context.Context, roomID string) ([]Element, error) {
	rows, err := d.pool.Query(ctx, 
		`SELECT id, room_id, type, world_coords, properties, created_at FROM elements WHERE room_id = $1`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var elements []Element
	for rows.Next() {
		var e Element
		if err := rows.Scan(&e.ID, &e.RoomID, &e.Type, &e.WorldCoords, &e.Properties, &e.CreatedAt); err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}
	return elements, nil
}
