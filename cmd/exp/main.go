package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
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
	fmt.Printf("Connected")

	// Create tables
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
	  id SERIAL PRIMARY KEY,
	  name TEXT,
	  email TEXT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS orders (
	  id SERIAL PRIMARY KEY,
	  user_id INT NOT NULL,
	  amount INT,
	  description TEXT
	);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")
}
