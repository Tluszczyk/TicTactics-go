package dynamo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"services/DatabaseService/database/types"
	databaseTypes "services/DatabaseService/database/types"
	"services/lib/log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
	log.Info("Started GetItemFromDatabase")

	tableName, key := input.TableName, input.Key

	log.Info("Marshall key")
	// Marshal key
	marshalled, err := d.MarshallDatabaseItem(&key)
	if err != nil {
		return nil, err
	}

	log.Info("Assert marshalled key is of type map[string]dynamoTypes.AttributeValue")
	// Assert marshalled key is of type map[string]dynamoTypes.AttributeValue
	av, ok := marshalled.(map[string]dynamoTypes.AttributeValue)
	if !ok {
		return nil, errors.New("marshalled key is not of type map[string]dynamoTypes.AttributeValue")
	}

	log.Info("Create get item input")
	// Create get item input
	dynamoInput := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       av,
	}

	log.Info("Execute get item")
	// Execute get item
	result, err := d.DynamodbClient.GetItem(context.Background(), dynamoInput)

	if err != nil {
		return nil, err
	}

	log.Info("Unmarshal result")
	// Unmarshal result
	var item types.DatabaseItem
	err = d.UnmarshallDatabaseItem(result.Item, &item)

	if err != nil {
		return nil, err
	}

	return &types.DatabaseGetItemOutput{
		Item: item,
	}, nil
}

func (d *DynamoDatabaseService) PutItemInDatabase(input *types.DatabasePutItemInput) (*types.DatabasePutItemOutput, error) {
	log.Info("Started PutItemInDatabase")

	tableName, item := input.TableName, input.Item

	// Marshal key
	marshalled, err := d.MarshallDatabaseItem(&item)
	if err != nil {
		return nil, err
	}

	// Assert marshalled key is of type map[string]dynamoTypes.AttributeValue
	av, ok := marshalled.(map[string]dynamoTypes.AttributeValue)
	if !ok {
		return nil, errors.New("marshalled key is not of type map[string]dynamoTypes.AttributeValue")
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
	log.Info("Started DeleteItemFromDatabase")
	return nil, errors.New("not implemented")
}

func (d *DynamoDatabaseService) UpdateItemInDatabase(input *types.DatabaseUpdateItemInput) (*types.DatabaseUpdateItemOutput, error) {
	log.Info("Started UpdateItemInDatabase")
	return nil, errors.New("not implemented")
}

func (d *DynamoDatabaseService) QueryDatabase(input *types.DatabaseQueryInput) (*types.DatabaseQueryOutput, error) {
	log.Info("Started QueryDatabase")
	return nil, errors.New("not implemented")
}

func (d *DynamoDatabaseService) MarshallDatabaseItem(item *databaseTypes.DatabaseItem) (interface{}, error) {
	log.Info("Started MarshallDatabaseItem")

	result, err := attributevalue.MarshalMap(item.Attributes)

	if err != nil {
		return nil, err
	}

	PK, SK := "", ""

	for key, value := range item.PK {
		PK += fmt.Sprintf("%s#%s", key, value)
	}

	for key, value := range item.SK {
		SK += fmt.Sprintf("%s#%s", key, value)
	}

	result["PK"] = &dynamoTypes.AttributeValueMemberS{Value: PK}
	result["SK"] = &dynamoTypes.AttributeValueMemberS{Value: SK}

	return result, nil
}

// UnmarshallDatabaseItem unmarshalls a database-specific item into a DatabaseItem
// It takes a map[string]dynamoTypes.AttributeValue, which is the type returned by dynamodb.GetItemOutput.Item
func (d *DynamoDatabaseService) UnmarshallDatabaseItem(item interface{}, result *types.DatabaseItem) error {
	log.Info("Started UnmarshallDatabaseItem")

	// TODO

	dynamoItem, ok := item.(map[string]dynamoTypes.AttributeValue)
	if !ok {
		return errors.New("item is not of type map[string]dynamoTypes.AttributeValue")
	}

	handleKey := func(key string, m *map[databaseTypes.FieldType]interface{}) {
		K := dynamoItem[key].(*dynamoTypes.AttributeValueMemberS).Value
		kParts := strings.Split(K, "#")

		for i := 0; i < len(kParts); i += 2 {
			key := databaseTypes.FieldType(kParts[i])
			value := kParts[i+1]

			(*m)[key] = value
		}
	}

	for key, value := range dynamoItem {
		switch key {
		case "PK":
			handleKey(key, &result.PK)
		case "SK":
			handleKey(key, &result.SK)
		default:
			result.Attributes[databaseTypes.FieldType(key)] = value // value is a pointer to an AttributeValue of type AttributeValueMemberS, AttributeValueMemberN, or AttributeValueMemberBOOL
		}
	}

	return nil
}
