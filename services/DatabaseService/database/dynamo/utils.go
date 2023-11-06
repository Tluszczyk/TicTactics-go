package dynamo

import (
	"fmt"
	databaseTypes "services/DatabaseService/database/types"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamoTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func MarshallDatabaseItem(item databaseTypes.DatabaseItem) (map[string]dynamoTypes.AttributeValue, error) {
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

func UnmarshallDatabaseItem(item map[string]dynamoTypes.AttributeValue, result *databaseTypes.DatabaseItem) error {

	// TODO

	handleKey := func(key string, m *map[databaseTypes.FieldType]interface{}) {
		K := item[key].(*dynamoTypes.AttributeValueMemberS).Value
		kParts := strings.Split(K, "#")

		for i := 0; i < len(kParts); i += 2 {
			key := databaseTypes.FieldType(kParts[i])
			value := kParts[i+1]

			(*m)[key] = value
		}
	}

	for key, value := range item {
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
