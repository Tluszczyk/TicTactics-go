package dynamo

import (
	"reflect"
	"testing"

	databaseTypes "services/DatabaseService/database/types"

	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestMarshallDatabaseItem(t *testing.T) {
	item := databaseTypes.DatabaseItem{
		Attributes: map[databaseTypes.FieldType]interface{}{
			"attribute1": "value1",
			"attribute2": 2,
			"attribute3": true,
		},
		PK: map[databaseTypes.FieldType]interface{}{
			databaseTypes.EMAIL: "pkValue",
		},
		SK: map[databaseTypes.FieldType]interface{}{
			databaseTypes.USERNAME: "skValue",
		},
	}

	expected := map[string]dynamoTypes.AttributeValue{
		"attribute1": &dynamoTypes.AttributeValueMemberS{Value: "value1"},
		"attribute2": &dynamoTypes.AttributeValueMemberN{Value: "2"},
		"attribute3": &dynamoTypes.AttributeValueMemberBOOL{Value: true},
		"PK":         &dynamoTypes.AttributeValueMemberS{Value: "EMAIL#pkValue"},
		"SK":         &dynamoTypes.AttributeValueMemberS{Value: "USER#skValue"},
	}

	result, err := MarshallDatabaseItem(item)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	keys := []string{
		"attribute1",
		"attribute2",
		"attribute3",
		"PK",
		"SK",
	}

	for _, key := range keys {
		if !reflect.DeepEqual(result[key], expected[key]) {
			t.Errorf("Unexpected result. Expected: %v, but got: %v", expected[key], result[key])
		}
	}
}

func TestUnmarshallDatabaseItem(t *testing.T) {
	item := map[string]dynamoTypes.AttributeValue{
		"attribute1": &dynamoTypes.AttributeValueMemberS{Value: "value1"},
		"attribute2": &dynamoTypes.AttributeValueMemberN{Value: "2"},
		"attribute3": &dynamoTypes.AttributeValueMemberBOOL{Value: true},
		"PK":         &dynamoTypes.AttributeValueMemberS{Value: "EMAIL#pkValue"},
		"SK":         &dynamoTypes.AttributeValueMemberS{Value: "USER#skValue"},
	}

	expected := databaseTypes.DatabaseItem{
		Attributes: map[databaseTypes.FieldType]interface{}{
			"attribute1": "value1",
			"attribute2": 2,
			"attribute3": true,
		},
		PK: map[databaseTypes.FieldType]interface{}{
			databaseTypes.EMAIL: "pkValue",
		},
		SK: map[databaseTypes.FieldType]interface{}{
			databaseTypes.USERNAME: "skValue",
		},
	}

	result := databaseTypes.NewDatabaseItem() // TODO
	err := UnmarshallDatabaseItem(item, result)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result.PK, expected.PK) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected.PK, result.PK)
	}

	if !reflect.DeepEqual(result.SK, expected.SK) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected.SK, result.SK)
	}

	if !reflect.DeepEqual(result.Attributes, expected.Attributes) {
		t.Errorf("Unexpected result. Expected: %v, but got: %v", expected.Attributes, result.Attributes)
	}
}
