package main

import (
	"errors"
	"reflect"
	"services/DatabaseService/database"
	databaseTypes "services/DatabaseService/database/types"
	messageTypes "services/lib/types"
	"testing"

	"github.com/google/uuid"
)

type mockDatabaseService struct {
	database.DatabaseService

	queryResult []databaseTypes.DatabaseItem
	err         error
}

func (m mockDatabaseService) GetItemFromDatabase(tableName string, itemID string) (databaseTypes.DatabaseItem, error) {
	return databaseTypes.DatabaseItem{}, nil
}

func (m mockDatabaseService) PutItemInDatabase(tableName string, itemID string, item databaseTypes.DatabaseItem) error {
	return m.err
}

func (m mockDatabaseService) DeleteItemFromDatabase(tableName string, itemID string) error {
	return m.err
}

func (m mockDatabaseService) UpdateItemInDatabase(tableName string, itemID string, item databaseTypes.DatabaseItem) error {
	return m.err
}

func (m mockDatabaseService) QueryDatabase(tableName string, query databaseTypes.DatabaseQuery) (databaseTypes.DatabaseQueryResult, error) {
	return m.queryResult, m.err
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
					"uid":      uuid.New().String(),
					"username": "testuser",
					"hash_id":  uuid.New().String(),
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
			databaseErr:    errors.New("database error"),
			expectedResult: false,
			expectedErr:    errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB := mockDatabaseService{
				queryResult: test.databaseResult,
				err:         test.databaseErr,
			}

			result, err := CheckIfUserAlreadyExists(mockDB, "users", "testuser")

			if result != test.expectedResult {
				t.Errorf("expected %v, but got %v", test.expectedResult, result)
			}

			if !reflect.DeepEqual(err, test.expectedErr) {
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
			databaseErr: errors.New("database error"),
			expectedErr: errors.New("database error"),
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

			if !reflect.DeepEqual(err, test.expectedErr) {
				t.Errorf("expected %v, but got %v", test.expectedErr, err)
			}
		})
	}
}
