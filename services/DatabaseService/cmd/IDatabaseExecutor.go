package main

import (
	lib "server/services/lib"
)

type IDatabaseExecutor interface {
	SaveUser(lib.Token, lib.User) lib.Status
	DeleteUser(lib.Token, lib.UserID) lib.Status
	UserExists(lib.Token, lib.UserID) (bool, lib.Status)

	SaveCredentials(lib.Token, lib.Credentials) lib.Status
	DeleteCredentials(lib.Token, lib.UserID) lib.Status
	CredentialsExist(lib.Token, lib.Credentials) (bool, lib.Status)

	SaveSession(lib.Token, lib.Session) lib.Status
	DeleteSession(lib.Token, lib.UserID) lib.Status
	GetSession(lib.Token, lib.UserID) (lib.Session, lib.Status)
}
