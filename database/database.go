package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"my_project/utils/password"
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

type User struct {
	Id             int       `json:"id,omitempty" xml:"id,omitempty"`
	Name           string    `json:"name,omitempty" xml:"name,omitempty"`
	Email          string    `json:"email,omitempty" xml:"email,omitempty"`
	Login          string    `json:"login,omitempty" xml:"login,omitempty"`
	HashedPassword string    `json:"hash,omitempty" xml:"hash,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty" xml:"created_at,omitempty"`
}

func CreateUser(u *User, db *pgx.Conn) error {
	var conflictField string

	err := db.QueryRow(
		context.Background(),
		`SELECT
        CASE 
            WHEN email = $1 THEN 'email'
            WHEN login = $2 THEN 'login'
        END as conflict_field
     	FROM users WHERE email = $1 OR login = $2`,
		u.Email, u.Login,
	).Scan(&conflictField)

	if err == nil {
		return fmt.Errorf("user with %s '%s' already exists", conflictField,
			map[string]string{"email": u.Email, "login": u.Login}[conflictField])
	}

	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	_, err = db.Exec(
		context.Background(),
		`INSERT INTO users (name, email, login, hashed_password)  
		VALUES ($1, $2, $3, $4)`,
		u.Name, u.Email, u.Login, u.HashedPassword)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return err
}

func ReadUser(id int, db *pgx.Conn) (*User, error) {
	var user User

	err := db.QueryRow(
		context.Background(),
		`SELECT id, name, email, login
		FROM users
		WHERE id = $1`,
		id).Scan(&user.Id, &user.Name, &user.Email, &user.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to read user: %w", err)
	}

	return &user, nil
}

func UpdateUser(id int, u *User, db *pgx.Conn) error {
	var (
		exists bool
		err    error
	)

	err = db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	err = db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)",
		u.Email, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	}
	if exists {
		return fmt.Errorf("email %s is already taken", u.Email)
	}

	err = db.QueryRow(
		context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 AND id != $2)",
		u.Login, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check login uniqueness: %w", err)
	}
	if exists {
		return fmt.Errorf("login %s is already taken", u.Login)
	}

	if u.HashedPassword != "" {
		_, err = db.Exec(
			context.Background(),
			`UPDATE users 
            SET name = $1, email = $2, login = $3, hashed_password = $4
            WHERE id = $5`,
			u.Name, u.Email, u.Login, u.HashedPassword, id)
	} else {
		_, err = db.Exec(
			context.Background(),
			`UPDATE users 
            SET name = $1, email = $2, login = $3
            WHERE id = $4`,
			u.Name, u.Email, u.Login, id)
	}

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return err
}

func DeleteUser(id int, db *pgx.Conn) error {
	var exists bool
	err := db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %d not found", id)
	}

	cT, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := cT.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no user was deleted")
	}

	return nil
}

func Authenticate(u User, db *pgx.Conn) (bool, error) {
	bdPerson := User{}
	err := db.QueryRow(
		context.Background(),
		`SELECT email, login, hashed_password 
		FROM users 
		WHERE email = $1 OR login = $1`, u.Login).Scan(&bdPerson.Email, &bdPerson.Login, &bdPerson.HashedPassword)
	if err != nil {
		return false, err
	}

	if (u.Login == bdPerson.Email || u.Login == bdPerson.Login) && password.Check(u.HashedPassword, bdPerson.HashedPassword) {
		return true, nil
	}

	return false, fmt.Errorf("invalid credentials")
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
		rows, err = conn.Query(context.Background(), `SELECT products.id, products.name, categories.name
			FROM products LEFT JOIN categories 
			ON products.category_id = categories.id`)
	} else {
		rows, err = conn.Query(context.Background(), `SELECT products.id, products.name, categories.name
			FROM products LEFT JOIN categories 
			ON products.category_id = categories.id 
			WHERE categories.name = $1`, category)
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
	row := conn.QueryRow(context.Background(), "SELECT name, description, price, amount, created_at FROM products WHERE id = $1", id)
	err := row.Scan(&p.Name, &p.Description, &p.Price, &p.Amount, &p.CreatedAt)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(p.Name, p.Description, p.Price, p.Amount, p.CreatedAt)

	return p
}
