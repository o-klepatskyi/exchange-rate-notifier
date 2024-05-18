package database

import (
    "database/sql"
    "fmt"
    "os"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)
    conn, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		os.Exit(1)
	} else {
		fmt.Println("Initialized database")
		db = conn
	}
}

func CreateTable() {
    query := `
    CREATE TABLE IF NOT EXISTS emails (
        email VARCHAR(255) PRIMARY KEY
    );`
    _, err := db.Exec(query)
    if err != nil {
        fmt.Println("Error creating table:", err)
        os.Exit(1)
    }
}

func AddEmail(email string) error {
	if db == nil {
		return fmt.Errorf("database is not initialized")
	}
	_, err := db.Exec("INSERT INTO emails (email) VALUES ($1)", email)
    return err
}

func GetAllEmails() ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}
	rows, err := db.Query("SELECT email FROM emails")
	if err != nil {
        fmt.Println("Error fetching emails:", err)
        return nil, err
    }
	defer rows.Close()

	var emails []string

    for rows.Next() {
        var email string
        if err := rows.Scan(&email); err != nil {
            return nil, err
        }
        emails = append(emails, email)
    }

    return emails, nil
}
