package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Users struct {
	id         int
	name       string
	email      string
	created_at time.Time
}

const DATABASE_URL_1 = "postgres://akg:helloAkg@localhost:5432/coffee_can?sslmode=disable"
const DATABASE_URL_2 = "postgres://akgReplica:helloDB2@localhost:5432/replica_coffee_can?sslmode=disable"

// Create database connection
func DBConn() *sql.DB {
	dbConn, err := sql.Open("pgx", DATABASE_URL_1)
	if err != nil {
		fmt.Println("Unable to connect to the DB: ", err)
	}

	err = dbConn.Ping()
	if err != nil {
		fmt.Println("Unable to ping to the DB: ", err)
	}
	return dbConn
}

// Add new user
func addUser(dbconn *sql.DB, name string, email string) {
	query := "INSERT INTO users (name, email) VALUES ($1, $2)"
	_ = dbconn.QueryRow(query, name, email)
}

// Update new user
func updateUser(dbconn *sql.DB, Id int, email string) {
	query := "UPDATE SET users email=$1 WHERE Id=$2"
	_ = dbconn.QueryRow(query, Id, email)
}

// Gell all user
func getAllUsers(dbconn *sql.DB) (*sql.Rows, error) {
	query := "SELECT * FROM users"
	rows, err := dbconn.Query(query)
	return rows, err
}

func main() {
	// Get database connection
	dbconn := DBConn()
	defer dbconn.Close()

	// Add new user
	addUser(dbconn, "John", "john@wwe.com")

	// Get all users
	rows, err := getAllUsers(dbconn)
	if err != nil {
		fmt.Println("Unable fetch records: ", err)
	}
	for rows.Next() {
		user := new(Users)

		err := rows.Scan(&user.id, &user.name, &user.email, &user.created_at)
		if err != nil {
			fmt.Println("Unable to read user: ", err)
		}
		fmt.Printf("Id: %d, name: %s, email: %s, created_at: %v\n", user.id, user.name, user.email, user.created_at)
		defer rows.Close()
	}

	// Update user - 1
	updateUser(dbconn, 1, "Robart@GOT.com")
}
