package database

import (
	"errors"
	"reflect"
	"services/DatabaseService/database/dynamo"
	"services/DatabaseService/database/types"
	"testing"
)

type mockDatabaseService struct {
	DatabaseService
}

func (m *mockDatabaseService) GetItemFromDatabase(tableName string, key string) (types.DatabaseItem, error) {
	return types.DatabaseItem{}, nil
}

func (m *mockDatabaseService) PutItemInDatabase(tableName string, key string, item types.DatabaseItem) error {
	return nil
}

func (m *mockDatabaseService) DeleteItemFromDatabase(tableName string, key string) error {
	return nil
}

func (m *mockDatabaseService) UpdateItemInDatabase(tableName string, key string, item types.DatabaseItem) error {
	return nil
}

func (m *mockDatabaseService) QueryDatabase(tableName string, query types.DatabaseQuery) (types.DatabaseQueryResult, error) {
	return types.DatabaseQueryResult{}, nil
}

func TestGetDatabaseService(t *testing.T) {
	t.Run("returns a DynamoDatabaseService when deployment option is DYNAMO", func(t *testing.T) {
		service, err := GetDatabaseService(DYNAMO)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		_, ok := service.(*dynamo.DynamoDatabaseService)
		if !ok {
			t.Errorf("expected a DynamoDatabaseService, got %T", service)
		}
	})

	t.Run("returns an error when deployment option is unsupported", func(t *testing.T) {
		_, err := GetDatabaseService(NONE)
		if !errors.Is(err, errors.ErrUnsupported) {
			t.Errorf("expected error %v, got %v", errors.ErrUnsupported, err)
		}
	})
}

func TestDatabaseServiceMethods(t *testing.T) {
	service := &mockDatabaseService{}

	t.Run("GetItemFromDatabase returns a DatabaseItem and no error", func(t *testing.T) {
		item, err := service.GetItemFromDatabase("table", "key")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(item, types.DatabaseItem{}) {
			t.Errorf("expected an empty DatabaseItem, got %v", item)
		}
	})

	t.Run("PutItemInDatabase returns no error", func(t *testing.T) {
		err := service.PutItemInDatabase("table", "key", types.DatabaseItem{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("DeleteItemFromDatabase returns no error", func(t *testing.T) {
		err := service.DeleteItemFromDatabase("table", "key")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("UpdateItemInDatabase returns no error", func(t *testing.T) {
		err := service.UpdateItemInDatabase("table", "key", types.DatabaseItem{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("QueryDatabase returns a DatabaseQueryResult and no error", func(t *testing.T) {
		result, err := service.QueryDatabase("table", types.DatabaseQuery{})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(result, types.DatabaseQueryResult{}) {
			t.Errorf("expected an empty DatabaseQueryResult, got %v", result)
		}
	})
}
