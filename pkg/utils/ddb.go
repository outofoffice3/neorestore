package utils

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Function to check if a given prefix is in the prefixes list
func ContainsPrefix(item map[string]types.AttributeValue, prefix string) bool {
	if item != nil {
		attr, exists := item["prefixes"]
		if exists {
			list, isList := attr.(*types.AttributeValueMemberL)
			if isList {
				for _, val := range list.Value {
					str, isString := val.(*types.AttributeValueMemberS)
					if isString && str.Value == prefix {
						return true
					}
				}
			}
		}
	}
	return false
}

// Utility function to create a list of DynamoDB string attributes
func CreateStringListAttribute(items []string) []types.AttributeValue {
	result := make([]types.AttributeValue, len(items))
	for i, item := range items {
		result[i] = &types.AttributeValueMemberS{Value: item}
	}
	return result
}

func StructToMap(input interface{}) map[string]types.AttributeValue {
	keys, err := attributevalue.MarshalMap(input)
	if err != nil {
		return nil
	}
	return keys
}
