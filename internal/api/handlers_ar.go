package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/AleksKAG/ar-backend/internal/store"
)

// GetElementsForAR Возвращает все элементы помещения с их 3D-координатами для AR-движка
func GetElementsForAR(db *store.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID := vars["room_id"]

		elements, err := db.GetElementsByRoomID(r.Context(), roomID)
		if err != nil {
			http.Error(w, "failed to fetch elements", http.StatusInternalServerError)
			return
		}

		// AR-клиент получит JSON вида: 
		// [{"id": "...", "type": "socket", "world_coords": {"x":1.2, "y":0.9, "z":0.0}, "properties": {...}}]
		json.NewEncoder(w).Encode(elements)
	}
}
