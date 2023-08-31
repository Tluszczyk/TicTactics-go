package lib

import (
	"github.com/google/uuid"
)

type Credentials struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type UserID uuid.UUID
type User struct {
	UserID   UserID
	Username string
}

type Token uuid.UUID

type Code int

const (
	OK Code = iota
	ERR
)

type Status struct {
	Code    Code
	Message string
}

type Session struct {
	Token Token
	User  User
}
