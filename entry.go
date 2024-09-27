package main

import (
	"log"
	"medods_jwt_service/logic"

	_ "github.com/lib/pq"
)

func main() {
	userDB, err := new(logic.UserDB).InitUserDB()
	if err != nil {
		log.Fatalf("Fatalno zadolbal err != nil")
	}
	defer userDB.Close()
}
