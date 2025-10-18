package user

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"my_project/database"
	"my_project/utils/password"
	"net/http"
	"strconv"
)

type Handler struct {
	db *pgx.Conn
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		user database.User
	)

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	user.HashedPassword, err = password.Hash(user.HashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateUser(&user, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"message":  "User created successfully",
		"redirect": "/login",
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		user database.User
	)

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Printf("JSON decode error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	if user.Name == "" || user.Email == "" || user.Login == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err = database.ReadUser(user.Id, h.db)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.HashedPassword != "" {
		user.HashedPassword, err = password.Hash(user.HashedPassword)
		if err != nil {
			http.Error(w, "Password hashing failed: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = database.UpdateUser(user.Id, &user, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User update successfully",
	})
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, `{"error": "User ID is required"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}
	user, err := database.ReadUser(id, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== DELETE PERSON CALLED ===")
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	user := database.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	access, err := database.Authenticate(user, h.db)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database error: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if access {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Login successful",
			"redirect": "/market",
		})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid email or password",
		})
	}
}
