package types

import (
	"testing"
)

func TestDatabaseQueryResult(t *testing.T) {
	// Test empty result
	var emptyResult DatabaseQueryResult
	if len(emptyResult) != 0 {
		t.Errorf("Expected empty result, but got %v", emptyResult)
	}

	// Test non-empty result
	result := DatabaseQueryResult{
		DatabaseItem{"id": 1, "name": "Alice"},
		DatabaseItem{"id": 2, "name": "Bob"},
	}
	if len(result) != 2 {
		t.Errorf("Expected result length of 2, but got %v", len(result))
	}
	if result[0]["id"] != 1 {
		t.Errorf("Expected first item id to be 1, but got %v", result[0]["id"])
	}
	if result[1]["name"] != "Bob" {
		t.Errorf("Expected second item name to be Bob, but got %v", result[1]["name"])
	}
}
