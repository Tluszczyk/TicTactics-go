package mongo

import (
	"context"
	"errors"
	"fmt"
	"services/DatabaseService/database/types"
	"services/lib/log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabaseService struct {
	MongodbClient    *mongo.Client
	database         *mongo.Database
	cancelConnection context.CancelFunc
}

func NewMongoDatabaseService() (*MongoDatabaseService, error) {

	// Get environment variables
	endpoint_url := "localhost:27017" // os.Getenv("ENDPOINT_URL")

	database_name := "TicTactics" // os.Getenv("DATABASE_NAME")

	// Load default config and assign region
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", endpoint_url)))
	if err != nil {
		log.Error("Failed to connect to MongoDB!")
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error("Failed to connect to MongoDB!")
		return nil, err
	}
	log.Info("Connected to MongoDB!")

	// Create and return AWS client
	return &MongoDatabaseService{
		MongodbClient:    client,
		database:         client.Database(database_name),
		cancelConnection: cancel,
	}, nil
}

func (d *MongoDatabaseService) Disconnect() error {
	d.cancelConnection()
	return d.MongodbClient.Disconnect(context.Background())
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

	// Convert marshalled key to a mongo filter item
	filter, err := bson.Marshal(marshalled)
	if err != nil {
		return nil, errors.New("marshalled key is not of type bson.D")
	}

	fmt.Printf("%+v\n", marshalled)

	// Execute get item
	var result bson.D
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Marshal result to BSON
	resultAsBSON, err := bson.Marshal(result)
	if err != nil {
		return nil, err
	}

	// Unmarshal result to map[string]interface{}
	var resultAsMap map[string]interface{}
	err = bson.Unmarshal(resultAsBSON, &resultAsMap)
	if err != nil {
		return nil, err
	}

	// Unmarshal result to DatabaseItem
	var item types.DatabaseItem
	err = d.UnmarshallDatabaseItem(resultAsMap, &item)
	if err != nil {
		return nil, err
	}

	// Return result
	return &types.DatabaseGetItemOutput{
		Item: item,
	}, nil
}

func (d *MongoDatabaseService) PutItemInDatabase(input *types.DatabasePutItemInput) (*types.DatabasePutItemOutput, error) {
	// Get the collection
	collectionName := input.TableName
	collection := d.database.Collection(collectionName)

	// Marshal item
	marshalled, err := d.MarshallDatabaseItem(&input.Item)
	if err != nil {
		return nil, err
	}

	// Execute put item
	_, err = collection.InsertOne(context.Background(), marshalled)
	if err != nil {
		return nil, err
	}

	// Return result
	return &types.DatabasePutItemOutput{}, nil
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

// MarshallDatabaseItem marshalls a DatabaseItem into a flat interface
func (d *MongoDatabaseService) MarshallDatabaseItem(item *types.DatabaseItem) (interface{}, error) {
	result := make(map[string]interface{})

	insert := func(prefix string, delimeter string, mapping map[types.FieldType]interface{}) {
		for key, value := range mapping {
			fieldKey := fmt.Sprintf("%s%s%s", prefix, delimeter, key)
			fieldValue := value

			result[fieldKey] = fieldValue
		}
	}

	insert("PK", "#", item.PK)
	insert("SK", "#", item.SK)
	insert("", "", item.Attributes)

	return result, nil
}

// UnmarshallDatabaseItem unmarshalls a flat map item into a DatabaseItem
func (d *MongoDatabaseService) UnmarshallDatabaseItem(item interface{}, result *types.DatabaseItem) error {
	// Initialize result
	result.PK = make(map[types.FieldType]interface{})
	result.SK = make(map[types.FieldType]interface{})
	result.Attributes = make(map[types.FieldType]interface{})

	// Assert marshalled item is of type map[string]interface{}
	av, ok := item.(map[string]interface{})
	if !ok {
		return errors.New("marshalled item is not of type map[string]interface{}")
	}

	// Unmarshal item
	for key, value := range av {
		if len(key) == 0 {
			continue
		}

		field := strings.Split(key, "#")

		var fieldKeyType string
		var fieldType string

		if len(field) == 1 {
			fieldKeyType = "Attributes"
			fieldType = field[0]
		} else {
			fieldKeyType = field[0]
			fieldType = field[1]
		}

		// Get field value
		fieldValue := value

		// Set field value
		switch fieldKeyType {
		case "PK":
			result.PK[types.FieldType(fieldType)] = fieldValue
		case "SK":
			result.SK[types.FieldType(fieldType)] = fieldValue
		case "Attributes":
			result.Attributes[types.FieldType(fieldType)] = fieldValue
		}
	}

	return nil
}
