package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"my_project/database"
	"net/http"
	"strings"
)

type Handler struct {
	db *pgx.Conn // или ваш тип подключения к БД
}

func NewHandler(db *pgx.Conn) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.Path)

	name := r.URL.Query().Get("name")
	if name == "" {
		return
	}

	fmt.Printf("Name: %s\n", name)

	p := database.GetPersonInfo(name, h.db)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.Path)

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
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.Path)

	path := r.URL.Path
	parts := strings.Split(path, "/")

	pr := database.GetProduct(parts[2], h.db)

	//json, err := json.Marshal(pr)
	//
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")
	////w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pr)

	//w.Write(json)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

	pr := database.Person{}

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(pr)
	err := database.CreateUser(&pr, h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//json, err := json.Marshal(pr)
	//
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//w.Write(json)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *Handler) CheckUser(w http.ResponseWriter, r *http.Request) {
	pr := database.Person{}

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	acces, err := database.CheckUser(&pr, h.db)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database error: " + err.Error(),
		})
		return // ← ДОБАВЬТЕ ЭТОТ return
	}

	w.Header().Set("Content-Type", "application/json")

	if acces {
		fmt.Println("Доступ разрешен")
		// УБЕРИТЕ http.Redirect - используйте JSON redirect
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  true,
			"message":  "Login successful",
			"redirect": "/market", // ← редирект в JSON
		})
	} else {
		fmt.Println("Доступ запрещен")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid email or password",
		})
	}
}

func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== UPDATE PERSON CALLED ===")
	fmt.Printf("Method: %s, URL: %s\n", r.Method, r.URL.Path)

	pr := database.Person{}

	// Читаем сырое тело для отладки
	//body, _ := io.ReadAll(r.Body)
	//fmt.Printf("Raw request body: %s\n", string(body))
	//
	//// Сбрасываем Body чтобы можно было декодировать снова
	//r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		fmt.Printf("JSON decode error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	fmt.Printf("Получены данные: ID=%d, Name=%s, Email=%s, Hash=%s\n",
		pr.Id, pr.Name, pr.Email, pr.Hash)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Данные обновлены",
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/register", http.StatusFound)
}

func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== DELETE PERSON CALLED ===")
}
