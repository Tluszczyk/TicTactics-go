package dynamo

import (
	"context"
	"errors"
	"os"
	"services/DatabaseService/database/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDatabaseService struct {
	DynamodbClient *dynamodb.Client
}

func NewDynamoDatabaseService() (*DynamoDatabaseService, error) {

	// Get region from environment variable
	region := os.Getenv("AWS_REGION")
	endpoint_url := os.Getenv("ENDPOINT_URL")

	// Load default config and assign region
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           endpoint_url,
						SigningRegion: region,
					}, nil
				},
			),
		),
	)

	if err != nil {
		return nil, err
	}

	// Create and return AWS client
	return &DynamoDatabaseService{
		DynamodbClient: dynamodb.NewFromConfig(cfg),
	}, nil
}

func (d *DynamoDatabaseService) GetItemFromDatabase(input *types.DatabaseGetItemInput) (*types.DatabaseGetItemOutput, error) {
	tableName, key := input.TableName, input.Key

	// Marshal key
	av, err := MarshallDatabaseItem(key)

	if err != nil {
		return nil, err
	}

	// Create get item input
	dynamoInput := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       av,
	}

	// Execute get item
	result, err := d.DynamodbClient.GetItem(context.Background(), dynamoInput)

	if err != nil {
		return nil, err
	}

	// Unmarshal result
	var item types.DatabaseItem
	err = UnmarshallDatabaseItem(result.Item, &item)

	if err != nil {
		return nil, err
	}

	return &types.DatabaseGetItemOutput{
		Item: item,
	}, nil
}

func (d *DynamoDatabaseService) PutItemInDatabase(input *types.DatabasePutItemInput) (*types.DatabasePutItemOutput, error) {
	tableName, item := input.TableName, input.Item

	// Marshal item
	av, err := attributevalue.MarshalMap(item)

	if err != nil {
		return nil, err
	}

	// Create put item input
	dynamoInput := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      av,
	}

	// Execute put item
	_, err = d.DynamodbClient.PutItem(context.Background(), dynamoInput)

	if err != nil {
		return nil, err
	}

	return &types.DatabasePutItemOutput{}, nil
}

func (d *DynamoDatabaseService) DeleteItemFromDatabase(input *types.DatabaseDeleteItemInput) (*types.DatabaseDeleteItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (d *DynamoDatabaseService) UpdateItemInDatabase(input *types.DatabaseUpdateItemInput) (*types.DatabaseUpdateItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (d *DynamoDatabaseService) QueryDatabase(input *types.DatabaseQueryInput) (*types.DatabaseQueryOutput, error) {
	return nil, errors.New("not implemented")
}
