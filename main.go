package main

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"my_project/database"
	"my_project/server/handler/html"
	"my_project/server/handler/market"
	"my_project/server/handler/user"
	"net/http"
)

var db *pgx.Conn

var idClient int

func main() {

	db = database.Connect()
	defer database.Close(db)

	u := user.NewHandler(db)
	m := market.NewHandler(db)

	// CRUD operation user
	http.HandleFunc("POST /register", u.Create)
	http.HandleFunc("GET /user/{id}", u.Read)
	http.HandleFunc("PUT /user/update", u.Update)
	http.HandleFunc("DELETE /user/delete", u.Delete)

	http.HandleFunc("POST /login", u.LoginUser)

	// HTML form pages
	http.HandleFunc("GET /user/update", html.UserUpdatePage)
	http.HandleFunc("GET /register", html.RegistrationPage)
	http.HandleFunc("GET /login", html.LoginPage)

	// CRUD operation market
	http.HandleFunc("GET /market", m.GetProducts)
	http.HandleFunc("GET /api/market/{id}", m.GetProduct)

	// Redirect
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	http.HandleFunc("GET /market/{id}", html.MarketProductPage)

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Println(err)
	}
}
