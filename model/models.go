package model

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name       string
	Email      string
	GUID       int
	IP         string
	Created_at time.Time
	Token      TokenPair
}

type TokenPair struct {
	JWT_token
	Refresh_token
}

type JWT_token struct {
	Value            string
	HashRefreshToken string
	SessionID        string
}

type Refresh_token struct {
	Value        string
	SessionID    string
	Hash         string
	Times_to_use int
}

type Payload struct {
	UserID    int
	Email     string
	SessionID string
	IP        string
	jwt.RegisteredClaims
}

type jwt_error struct {
	err_type string
	message  string
}

func (err *jwt_error) Error() string {
	return fmt.Sprintf("Err type %s. Message: %s", err.err_type, err.message)
}

func NewError(err_type, message string) error {
	return &jwt_error{err_type: err_type, message: message}
}
