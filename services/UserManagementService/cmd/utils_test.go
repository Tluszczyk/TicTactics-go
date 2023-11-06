package main

import (
	"errors"
	"services/DatabaseService/database"
	databaseTypes "services/DatabaseService/database/types"
	messageTypes "services/lib/types"
	"testing"
)

type mockDatabaseService struct {
	database.DatabaseService

	queryResult []databaseTypes.DatabaseItem
	err         error
}

func (m mockDatabaseService) GetItemFromDatabase(input *databaseTypes.DatabaseGetItemInput) (*databaseTypes.DatabaseGetItemOutput, error) {
	var item databaseTypes.DatabaseItem

	if m.queryResult == nil {
		return nil, m.err
	}

	if len(m.queryResult) > 0 {
		item = m.queryResult[0]
	} else {
		item = databaseTypes.DatabaseItem{}
	}

	return &databaseTypes.DatabaseGetItemOutput{
		Item: item,
	}, nil
}

func (m mockDatabaseService) PutItemInDatabase(input *databaseTypes.DatabasePutItemInput) (*databaseTypes.DatabasePutItemOutput, error) {
	return &databaseTypes.DatabasePutItemOutput{}, m.err
}

func (m mockDatabaseService) DeleteItemFromDatabase(input *databaseTypes.DatabaseDeleteItemInput) (*databaseTypes.DatabaseDeleteItemOutput, error) {
	return &databaseTypes.DatabaseDeleteItemOutput{}, m.err
}

func (m mockDatabaseService) UpdateItemInDatabase(input *databaseTypes.DatabaseUpdateItemInput) (*databaseTypes.DatabaseUpdateItemOutput, error) {
	return &databaseTypes.DatabaseUpdateItemOutput{}, m.err
}

func (m mockDatabaseService) QueryDatabase(input *databaseTypes.DatabaseQueryInput) (*databaseTypes.DatabaseQueryOutput, error) {
	return &databaseTypes.DatabaseQueryOutput{
		Items: m.queryResult,
	}, m.err
}

type databaseError struct{}

func (d databaseError) Error() string {
	return "database error"
}

func TestCheckIfUserAlreadyExists(t *testing.T) {
	tests := []struct {
		name           string
		databaseResult []databaseTypes.DatabaseItem
		databaseErr    error
		expectedResult bool
		expectedErr    error
	}{
		{
			name: "user exists",
			databaseResult: []databaseTypes.DatabaseItem{
				{
					PK:         map[databaseTypes.FieldType]interface{}{databaseTypes.EMAIL: "test@test.com"},
					SK:         map[databaseTypes.FieldType]interface{}{databaseTypes.HASH_ID: "123546"},
					Attributes: map[databaseTypes.FieldType]interface{}{databaseTypes.PASSWORD_HASH: "testhash"},
				},
			},
			databaseErr:    nil,
			expectedResult: true,
			expectedErr:    nil,
		},
		{
			name:           "user does not exist",
			databaseResult: []databaseTypes.DatabaseItem{},
			databaseErr:    nil,
			expectedResult: false,
			expectedErr:    nil,
		},
		{
			name:           "database error",
			databaseResult: nil,
			databaseErr:    databaseError{},
			expectedResult: false,
			expectedErr:    databaseError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB := mockDatabaseService{
				queryResult: test.databaseResult,
				err:         test.databaseErr,
			}

			result, err := CheckIfUserAlreadyExists(mockDB, "users", messageTypes.Credentials{
				Username:     "testuser",
				Email:        "test@test.com",
				PasswordHash: "testhash",
			})

			if result != test.expectedResult {
				t.Errorf("expected %v, but got %v", test.expectedResult, result)
			}

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("expected %v, but got %v", test.expectedErr, err)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		databaseErr error
		expectedErr error
	}{
		{
			name:        "success",
			databaseErr: nil,
			expectedErr: nil,
		},
		{
			name:        "database error",
			databaseErr: databaseError{},
			expectedErr: databaseError{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB := mockDatabaseService{
				err: test.databaseErr,
			}

			credentials := messageTypes.Credentials{
				Username:     "testuser",
				Email:        "testuser@test.com",
				PasswordHash: "testhash",
			}

			err := CreateUser(mockDB, "users", "password_hashes", credentials)

			if !errors.Is(err, test.expectedErr) {
				t.Errorf("expected %v, but got %v", test.expectedErr, err)
			}
		})
	}
}
