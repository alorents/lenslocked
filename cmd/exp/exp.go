package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/alorents/lenslocked/models"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func main() {
	postgresConfig := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		DBName:   "lenslocked",
		SSLMode:  "disable",
	}
	db, err := sql.Open("pgx", postgresConfig.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected\n")

	us := models.UserService{
		DB: db,
	}

	// Create tables
	//_, err = db.Exec(`
	//	CREATE TABLE users (
	//	id SERIAL PRIMARY KEY,
	//	email TEXT UNIQUE NOT NULL,
	//	password_hash TEXT NOT NULL
	//);
	//
	//CREATE TABLE IF NOT EXISTS orders (
	//  id SERIAL PRIMARY KEY,
	//  user_id INT NOT NULL,
	//  amount INT,
	//  description TEXT
	//);`)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Tables created.")

	// Create fake users
	startIdx := 0
	for i := startIdx; i <= startIdx+5; i++ {
		email := fmt.Sprintf("email_%d", i)
		password := fmt.Sprintf("password_%d", i)
		user, err := us.Create(email, password)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Created user: %v\n", *user)
	}

	//userID := id
	//for i := 1; i <= 5; i++ {
	//	amount := i * 100
	//	desc := fmt.Sprintf("Fake order #%d", i)
	//	_, err := db.Exec(`
	//INSERT INTO orders(user_id, amount, description)
	//VALUES($1, $2, $3)`, userID, amount, desc)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//fmt.Println("Created fake orders.")
	//
	//type Order struct {
	//	ID          int
	//	UserID      int
	//	Amount      int
	//	Description string
	//}
	//var orders []Order
	//
	//userID = 1 // Use the same ID you used in the previous lesson
	//rows, err := db.Query(`
	//	SELECT id, amount, description
	//	FROM orders
	//	WHERE user_id=$1`, userID)
	//if err != nil {
	//	panic(err)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var order Order
	//	order.UserID = userID
	//	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	//	if err != nil {
	//		panic(err)
	//	}
	//	orders = append(orders, order)
	//}
	//err = rows.Err()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Orders:", orders)
}
