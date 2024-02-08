package cmd

import (
	types "services/lib/types"
)

type CreateGameRequest struct {
	Session types.Session `json:"session"`
	Settings types.GameSettings `json:"settings"`
}

type CreateGameResponse struct {
	Status types.Status `json:"status"`
}