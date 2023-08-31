package main

import (
	lib "server/services/lib"

	"github.com/google/uuid"
)

type AuthorisationBaseExecutor struct {
	IAuthorisationExecutor
}

func Login(executor AuthorisationBaseExecutor, credentials lib.Credentials) (lib.Token, lib.Status) {
	var token lib.Token
	var status lib.Status

	return token, status
}

func Logout(executor AuthorisationBaseExecutor, token lib.Token, uid uuid.UUID) lib.Status {
	var status lib.Status

	return status
}

func Register(executor AuthorisationBaseExecutor, credentials lib.Credentials) (lib.Token, lib.Status) {
	var token lib.Token
	var status lib.Status

	return token, status
}

func Validate(executor AuthorisationBaseExecutor, token lib.Token, uid uuid.UUID) lib.Status {
	var status lib.Status

	return status
}
