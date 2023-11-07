package mongo

import (
	"context"
	"errors"
	"fmt"
	"services/DatabaseService/database/types"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseService struct {
	MongodbClient *mongo.Client
	database      *mongo.Database
}

func NewMongoDatabaseService() (*MongoDatabaseService, error) {

	// Get environment variables
	endpoint_url := "localhost:27017" // os.Getenv("ENDPOINT_URL")

	database_name := "TicTactics" // os.Getenv("DATABASE_NAME")

	// Load default config and assign region
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", endpoint_url)))

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Create and return AWS client
	return &MongoDatabaseService{
		MongodbClient: client,
		database:      client.Database(database_name),
	}, nil
}

func (d *MongoDatabaseService) GetItemFromDatabase(input *types.DatabaseGetItemInput) (*types.DatabaseGetItemOutput, error) {
	// Get the collection
	collectionName, key := input.TableName, input.Key
	collection := d.database.Collection(collectionName)

	// Marshal key
	marshalled, err := d.MarshallDatabaseItem(&key)
	if err != nil {
		return nil, err
	}

	// Assert marshalled key to a mongo filter item
	filter, ok := marshalled.(bson.D)
	if !ok {
		return nil, errors.New("marshalled key is not of type bson.D")
	}

	// Execute get item
	var result bson.D
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Unmarshal result
	var item types.DatabaseItem
	err = d.UnmarshallDatabaseItem(result, &item)
	if err != nil {
		return nil, err
	}

	// Return result
	return &types.DatabaseGetItemOutput{
		Item: item,
	}, nil
}

func (d *MongoDatabaseService) PutItemInDatabase(input *types.DatabasePutItemInput) (*types.DatabasePutItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (d *MongoDatabaseService) DeleteItemFromDatabase(input *types.DatabaseDeleteItemInput) (*types.DatabaseDeleteItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (d *MongoDatabaseService) UpdateItemInDatabase(input *types.DatabaseUpdateItemInput) (*types.DatabaseUpdateItemOutput, error) {
	return nil, errors.New("not implemented")
}

func (d *MongoDatabaseService) QueryDatabase(input *types.DatabaseQueryInput) (*types.DatabaseQueryOutput, error) {
	return nil, errors.New("not implemented")
}

// MarshallDatabaseItem marshalls a DatabaseItem into a bson.D mongo filter
func (d *MongoDatabaseService) MarshallDatabaseItem(item *types.DatabaseItem) (interface{}, error) {
	result := bson.D{}

	fieldNames := []string{"PK", "SK", "Attributes"}
	for i, mapping := range []map[types.FieldType]interface{}{
		item.PK,
		item.SK,
		item.Attributes,
	} {
		for key, value := range mapping {
			result = append(result, bson.E{
				Key:   fmt.Sprintf("%s#%s", fieldNames[i], key),
				Value: value,
			})
		}
	}

	return result, nil
}

// UnmarshallDatabaseItem unmarshalls a database-specific item into a DatabaseItem
func (d *MongoDatabaseService) UnmarshallDatabaseItem(item interface{}, result *types.DatabaseItem) error {
	// Assert marshalled item is of type bson.D
	av, ok := item.(bson.D)
	if !ok {
		return errors.New("marshalled item is not of type bson.D")
	}

	// Unmarshal item
	for _, value := range av {
		if len(value.Key) == 0 {
			continue
		}

		field := strings.Split(value.Key, "#")

		var fieldName string
		var fieldKey string

		if len(field) == 1 {
			fieldName = "Attributes"
			fieldKey = field[0]
		} else {
			fieldName = field[0]
			fieldKey = field[1]
		}

		// Get field value
		fieldValue := value.Value

		// Set field value
		switch fieldName {
		case "PK":
			result.PK[types.FieldType(fieldKey)] = fieldValue
		case "SK":
			result.SK[types.FieldType(fieldKey)] = fieldValue
		case "Attributes":
			result.Attributes[types.FieldType(fieldKey)] = fieldValue
		}
	}

	return nil
}
