package market

import (
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"my_project/database"
	"net/http"
)

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	pr := database.GetProducts(category, h.db)

	//json, err := json.Marshal(pr)
	//
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pr)

	//w.Write(json)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, `{"error": "User ID is required"}`, http.StatusBadRequest)
		return
	}

	product := database.GetProduct(idStr, h.db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)

}
