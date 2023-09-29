package main

import (
	types "services/lib/types"
)

type CreateUserRequest struct {
	Credentials types.Credentials `json:"credentials"`
}

type CreateUserResponse struct {
	Status types.Status `json:"status"`
}
