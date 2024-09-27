package tests

import (
	"database/sql"
	logic "medods_jwt_service/logic"
	"testing"

	_ "github.com/lib/pq"
)

func TestDBConnection(t *testing.T) {
	connStr := "user=postgres dbname=medods password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to open the database: %v", err)
	}
	defer db.Close()

	// Check if the connection is valid
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping the database: %v", err)
	}

	t.Log("Successfully connected to the database!")
}

func TestDBGetUserById(t *testing.T) {
	userDB, err := new(logic.UserDB).InitUserDB()
	if err != nil {
		t.Fatalf("Fatalno zadolbal err != nil")
	}
	defer userDB.Close()
	user, err := logic.GetUserById(2, userDB.GetDBConn())
	if err != nil {
		t.Fatalf("Failed to get user by id: %s", err)
	}
	if user.Name == "maxim" {
		t.Log("Test passed!")
	}
}

func TestDBAddUserIntoTable(t *testing.T) {
	userDB, err := new(logic.UserDB).InitUserDB()
	if err != nil {
		t.Fatalf("Fatalno zadolbal err != nil")
	}
	defer userDB.Close()
	err = logic.AddUserIntoTable(userDB.GetDBConn(), "pashka", "228.228.227.012", "pashka@mail.ru")
	if err != nil {
		t.Fatalf("Failed to add user: %s", err)
	}
	t.Log("Test passed!")
}

func TestDBDeleteUserFromTable(t *testing.T) {
	userDB, err := new(logic.UserDB).InitUserDB()
	if err != nil {
		t.Fatalf("Fatalno zadolbal err != nil")
	}
	defer userDB.Close()
	err = logic.DeleteUserFromTable(userDB.GetDBConn(), "name", "pashka")
	if err != nil {
		t.Fatalf("Failed to delete user from table: %s", err)
	}
	t.Log("Test passed!")
}
