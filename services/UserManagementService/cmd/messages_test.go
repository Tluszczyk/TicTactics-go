package cmd

import (
	"encoding/json"
	"reflect"
	"testing"

	types "services/lib/types"
)

func TestCreateUserRequest(t *testing.T) {
	// create test data
	creds := types.Credentials{
		Username:     "testuser",
		Email:        "testuser@test.com",
		PasswordHash: "testpasswordhash",
	}
	req := CreateUserRequest{
		Credentials: creds,
	}

	// marshal to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Errorf("Error marshaling CreateUserRequest to JSON: %v", err)
	}

	// unmarshal back to struct
	var unmarshaledReq CreateUserRequest
	err = json.Unmarshal(jsonData, &unmarshaledReq)
	if err != nil {
		t.Errorf("Error unmarshaling JSON to CreateUserRequest: %v", err)
	}

	// compare original and unmarshaled structs
	if !reflect.DeepEqual(req, unmarshaledReq) {
		t.Errorf("CreateUserRequest structs not equal. Expected: %v, got: %v", req, unmarshaledReq)
	}
}

func TestCreateUserResponse(t *testing.T) {
	// create test data
	status := types.Status{
		Code:    200,
		Message: "OK",
	}
	resp := CreateUserResponse{
		Status: status,
	}

	// marshal to JSON
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Errorf("Error marshaling CreateUserResponse to JSON: %v", err)
	}

	// unmarshal back to struct
	var unmarshaledResp CreateUserResponse
	err = json.Unmarshal(jsonData, &unmarshaledResp)
	if err != nil {
		t.Errorf("Error unmarshaling JSON to CreateUserResponse: %v", err)
	}

	// compare original and unmarshaled structs
	if !reflect.DeepEqual(resp, unmarshaledResp) {
		t.Errorf("CreateUserResponse structs not equal. Expected: %v, got: %v", resp, unmarshaledResp)
	}
}
