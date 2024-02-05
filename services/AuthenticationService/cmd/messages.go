package cmd

import (
	types "services/lib/types"
)

type ValidateSessionRequest struct {
	Session types.Session `json:"session"`
}

type ValidateSessionResponse struct {
	Status types.Status `json:"status"`
}

type CreateSessionRequest struct {
	Credentials string `json:"credentials"`
}

type CreateSessionResponse struct {
	Status types.Status `json:"status"`
	Session types.Session `json:"session"`
}

type DeleteSessionRequest struct {
	Session types.Session `json:"session"`
}

type DeleteSessionResponse struct {
	Status types.Status `json:"status"`
}