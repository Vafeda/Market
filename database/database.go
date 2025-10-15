package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"
)

func Connect() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:0411@localhost:5432/market?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func Close(conn *pgx.Conn) error {
	return conn.Close(context.Background())
}

type Person struct {
	Id        int       `json:"id,omitempty" xml:"id,omitempty"`
	Name      string    `json:"name,omitempty" xml:"name,omitempty"`
	Email     string    `json:"email,omitempty" xml:"email,omitempty"`
	Hash      string    `json:"hash,omitempty" xml:"hash,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" xml:"created_at,omitempty"`
}

func (p Person) String() string {
	return fmt.Sprintf("Name: %s, Email: %s, Hash: %s", p.Name, p.Email, p.Hash)
}

func GetPersonInfo(name string, conn *pgx.Conn) Person {

	p := Person{}

	row := conn.QueryRow(context.Background(), "SELECT id, name, email, password_hash FROM person WHERE name = $1", name)

	err := row.Scan(&p.Id, &p.Name, &p.Email, &p.Hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query row: %v\n", err)
	}

	return p
}

type Product struct {
	Id          int       `json:"id,omitempty" xml:"id,omitempty"`
	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
	Description string    `json:"description,omitempty" xml:"description,omitempty"`
	Price       float64   `json:"price,omitempty" xml:"price,omitempty"`
	Amount      int       `json:"amount,omitempty" xml:"amount,omitempty"`
	Category    string    `json:"category,omitempty" xml:"category,omitempty"`
	CreatedAt   time.Time `json:"created_At,omitempty" xml:"created_At,omitempty"`
}

func GetProducts(category string, conn *pgx.Conn) []Product {

	var rows pgx.Rows
	var err error

	if len(category) == 0 {
		rows, err = conn.Query(context.Background(), `SELECT product.id, product.name, category.name
			FROM product JOIN category 
			ON product.category_id = category.id`)
	} else {
		rows, err = conn.Query(context.Background(), `SELECT product.id, product.name, category.name
			FROM product JOIN category 
			ON product.category_id = category.id 
			WHERE category.name = $1`, category)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query row: %v\n", err)
	}
	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var product Product
		rows.Scan(&product.Id, &product.Name, &product.Category)

		fmt.Println(product.Id, product.Name, product.Category)
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to query row: %v\n", err)
	}

	return products
}

func GetProduct(id string, conn *pgx.Conn) Product {
	p := Product{}
	row := conn.QueryRow(context.Background(), "SELECT name, description, price, amount, created_at FROM product WHERE id = $1", id)
	err := row.Scan(&p.Name, &p.Description, &p.Price, &p.Amount, &p.CreatedAt)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p.Name, p.Description, p.Price, p.Amount, p.CreatedAt)

	return p
}

func CreateUser(p *Person, conn *pgx.Conn) error {
	var existingID int
	err := conn.QueryRow(context.Background(), "SELECT id FROM person WHERE email = $1", p.Email).Scan(&existingID)

	if err == nil {
		return fmt.Errorf("user with email %s already exists", p.Email)
	}

	if err != pgx.ErrNoRows {
		return fmt.Errorf("database error: %v", err)
	}

	pg, err := conn.Exec(context.Background(), `INSERT INTO person (name, email, password_hash) VALUES ($1, $2, $3)`, p.Name, p.Email, p.Hash)
	if err != nil {
		fmt.Println(err)
		return errors.New("Логин уже существует в базе данных")
	}
	fmt.Println(pg)

	return nil
}

func CheckUser(p *Person, conn *pgx.Conn) (bool, error) {
	bdPerson := Person{}
	row := conn.QueryRow(context.Background(), "SELECT name, email, password_hash FROM person WHERE email = $1", p.Email)
	err := row.Scan(&bdPerson.Name, &bdPerson.Email, &bdPerson.Hash)
	if err != nil {
		return false, err
	}

	if bdPerson.Email == p.Email && bdPerson.Hash == p.Hash {
		return true, nil
	}

	return false, errors.New("Acces denied")
}

func UpdateUser(p *Person, conn *pgx.Conn) error {
	row := conn.Exec(context.Background(), "UPDATE person")
}
