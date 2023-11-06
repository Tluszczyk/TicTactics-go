package dynamo

import (
	"fmt"
	"os"
	"services/DatabaseService/database/types"
	"testing"
)

var testService *DynamoDatabaseService

func init() {
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("ENDPOINT_URL", "http://localhost:4566")

	var err error
	testService, err = NewDynamoDatabaseService()

	if err != nil {
		panic(err)
	}
}

func TestGetItemFromDatabase(t *testing.T) {
	t.Skip("This test should be run manually on created test stack")

	input := types.DatabaseGetItemInput{
		TableName: "TicTactics",
		Key: types.DatabaseItem{
			PK: map[types.FieldType]interface{}{types.UID: "testuser"},
			SK: map[types.FieldType]interface{}{types.USERNAME: "testpass"},
		},
	}

	item, err := testService.GetItemFromDatabase(&input)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Get item from database: Success: %+v\n", item)
	}
}

func TestPutItemInDatabase(t *testing.T) {
	t.Skip("This test should be run manually on created test stack")

	putInput := types.DatabasePutItemInput{
		TableName: "TicTactics",
		Item: types.DatabaseItem{
			PK:         map[types.FieldType]interface{}{types.UID: "testuser"},
			SK:         map[types.FieldType]interface{}{types.USERNAME: "testpass"},
			Attributes: map[types.FieldType]interface{}{types.ELO: 1000},
		},
	}

	_, err := testService.PutItemInDatabase(&putInput)

	if err != nil {
		t.Error(err)
	}

	fmt.Print("Put item in database: Success\n")

	getInput := types.DatabaseGetItemInput{
		TableName: "TicTactics",
		Key: types.DatabaseItem{
			PK: map[types.FieldType]interface{}{types.UID: "testuser"},
			SK: map[types.FieldType]interface{}{types.USERNAME: "testpass"},
		},
	}

	getOutput, err := testService.GetItemFromDatabase(&getInput)

	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("Get item from database: Success: %+v\n", getOutput.Item)
	}
}
