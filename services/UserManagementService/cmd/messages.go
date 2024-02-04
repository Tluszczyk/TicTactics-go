package cmd

import (
	types "services/lib/types"
)

type CreateUserRequest struct {
	Credentials types.Credentials `json:"credentials"`
}

type CreateUserResponse struct {
	Status types.Status `json:"status"`
}

type GetUserRequest struct {
	Username string `json:"username"`
}

type GetUserResponse struct {
	Status types.Status `json:"status"`
	User   types.User   `json:"user"`
}