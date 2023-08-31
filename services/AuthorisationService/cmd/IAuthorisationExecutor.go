package main

import (
	lib "server/services/lib"

	"github.com/google/uuid"
)

type IAuthorisationExecutor interface {
	Login(lib.Credentials) (lib.Token, lib.Status)

	Logout(lib.Token, uuid.UUID) lib.Status

	Register(lib.Credentials) (lib.Token, lib.Status)

	Validate(lib.Token, uuid.UUID) lib.Status
}
