package dynamo

import (
	"context"
	"os"
	"services/DatabaseService/database/types"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDatabaseService struct {
	DynamodbClient *dynamodb.Client
}

func NewDynamoDatabaseService() (*DynamoDatabaseService, error) {

	// Get region from environment variable
	region := os.Getenv("AWS_REGION")

	// Load default config and assign region
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

	if err != nil {
		return nil, err
	}

	// Create and return AWS client
	return &DynamoDatabaseService{
		DynamodbClient: dynamodb.NewFromConfig(cfg),
	}, nil
}

func (d *DynamoDatabaseService) GetItemFromDatabase(tableName string, key string) (types.DatabaseItem, error) {
	return nil, nil
}

func (d *DynamoDatabaseService) PutItemInDatabase(tableName string, key string, item types.DatabaseItem) error {
	// // Create put item input
	// input := &dynamodb.PutItemInput{
	// 	TableName: &tableName,
	// 	Item:      dynamoattribute.MarshalMap(item),
	// }

	return nil
}

func (d *DynamoDatabaseService) DeleteItemFromDatabase(tableName string, key string) error {
	return nil
}

func (d *DynamoDatabaseService) UpdateItemInDatabase(tableName string, key string, item types.DatabaseItem) error {
	return nil
}

func (d *DynamoDatabaseService) QueryDatabase(tableName string, query types.DatabaseQuery) (types.DatabaseQueryResult, error) {
	// // Create key condition expression
	// var queryFilter

	// for key, value := range query {
	// 	keyConditionBuilder = expression.Name(key).Equal(expression.Value(value))

	// // Create query input
	// input := &dynamodb.QueryInput{
	// 	TableName: &tableName,
	// 	KeyConditionExpression: &query.KeyConditionExpression,
	// }

	// // Execute query
	// result, err := d.DynamodbClient.Query(context.Background(), input)

	// if err != nil {
	// 	return nil, err
	// }

	// // Create query result
	// var queryResult types.DatabaseQueryResult
	// err = dynamoattribute.UnmarshalListOfMaps(result.Items, &queryResult)

	// if err != nil {
	// 	return nil, err
	// }

	// return queryResult, nil
	return nil, nil
}
