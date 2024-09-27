package logic

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var connDB *sql.DB

func InitConnection() (*sql.DB, error) {
	connStr := "user=postgres dbname=medods password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	connDB = db
	return db, nil
}

func GetDBConn() (*sql.DB, error) {
	if connDB == nil {
		db, err := InitConnection()
		if err != nil {
			return nil, err
		}
		connDB = db
	}
	return connDB, nil
}

func CloseConnection(db *sql.DB) {
	db.Close()
}

func AddUserIntoTable(db *sql.DB, name, IP, email string) error {
	var id_int int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&id_int)
	if err != nil {
		return err
	}
	UID := uuid.New()
	created := time.Now()
	id_int += 1
	query := `
	INSERT INTO users (name, ip, id, email, created_at, id_int) 
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, created_at`
	err = db.QueryRow(query, name, IP, UID, email, created, id_int).Scan(&UID, &created)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQuery(param string, value any) string {
	return fmt.Sprintf("DELETE FROM users WHERE %s = '%v'", param, value)
}

func DeleteUserFromTable(db *sql.DB, param string, value any) error {
	query := DeleteQuery(param, value)
	_, err := db.Exec(query)
	return err
}
