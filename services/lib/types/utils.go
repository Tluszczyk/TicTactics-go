package types

import (
	"encoding/json"
)

func ParseRequestBody(requestBody string, target interface{}) error {
	err := json.Unmarshal([]byte(requestBody), &target)
	if err != nil {
		return err
	}

	return nil
}
