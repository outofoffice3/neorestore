package manifest

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/common/vault"
	"github.com/outofoffice3/neorestore/pkg/constants"
	registerTypes "github.com/outofoffice3/neorestore/pkg/types"
	"github.com/outofoffice3/neorestore/pkg/utils"
)

type Manifest interface {
	// get a prefix item from table
	GetPrefixItem(key interface{}) (map[string]types.AttributeValue, bool)
	// put a prefix item to table
	PutPrefixItem(item interface{}) error
	// update item date restored
	UpdateItemDateRestored(item interface{}) error
	// remove prefix item from table
	RemovePrefixItem(key interface{}) error
	// put prefix list
	PutPrefixList(string) error
	// get reserved prefix item from table
	GetPrefixList() map[string]types.AttributeValue
	// put reserved prefix item to table
	AddPrefixToPrefixList(prefix string) error
	// get table name
	GetTableName() string
	// query manifest
	Query(input *dynamodb.QueryInput) *dynamodb.QueryPaginator
}

var (
	sos             logger.Logger
	prefixListCache []string
	localSecrets    map[string]string
	cachedManifest  Manifest
)

type _Manifest struct {
	client    *dynamodb.Client
	tableName string
}

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	// get secrets from vault
	vault, err := vault.NewVault()
	// return errors
	if err != nil {
		panic(err)
	}
	// get secrets for dynamodb
	result, err := vault.GetSecret(constants.NeoDDBSecretName)
	// return errors
	if err != nil {
		panic(err)
	}
	sos.Debugf("secrets: [%+v]", result)
	// type assert that value is equal to string
	secrets, ok := result.(string)
	// return errors
	if !ok {
		panic(errors.New("type assert to string error"))
	}
	sos.Debugf("secrets coverted to string: [%+v]", secrets)
	// convert json string of secrets to a map
	// load values to local secrets map
	err = json.Unmarshal([]byte(secrets), &localSecrets)
	// return errors
	if err != nil {
		panic(err)
	}
	// create new manifest instance
	cachedManifest = NewManifest()
	/// get prefix list
	pl := cachedManifest.GetPrefixList()
	// unmarshall into []string and save to preflix list cache []string
	err = attributevalue.UnmarshalMap(pl, &prefixListCache)
	// return errors
	if err != nil {
		panic(err)
	}
	sos.Infof("manifest init completed")
}

// create new Manifest
func NewManifest() Manifest {
	// if cached manifest is not nil, return it
	if cachedManifest != nil {
		return cachedManifest
	}
	// if not, create new instance
	cfg, err := config.LoadDefaultConfig(context.Background())
	cfg.Region = "us-east-1"
	if err != nil {
		panic(err)
	}
	// get table name, pk and sk from secrets
	tableName := localSecrets["manifestTableName"]
	// create new instance and saved to cache
	cachedManifest = &_Manifest{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
	// return cachedManifest
	return cachedManifest
}

// query manifest
func (m *_Manifest) Query(input *dynamodb.QueryInput) *dynamodb.QueryPaginator {
	return dynamodb.NewQueryPaginator(m.client, input)
}

// update item date restored
func (m *_Manifest) UpdateItemDateRestored(item interface{}) error {
	// type assert  item is prefix item
	prefixItem, ok := item.(registerTypes.PrefixItem)
	// return errors
	if !ok {
		sos.Errorf("type assert to prefix item error")
		return nil
	}
	sos.Debugf("prefix item: [%+v]", prefixItem)
	// convert keys to map[string]types.attributevalue
	keys := utils.StructToMap(prefixItem.Keys)
	sos.Debugf("keys: [%+v]", keys)
	// update item date restored
	_, err := m.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName:        aws.String(m.tableName),
		Key:              keys,
		UpdateExpression: aws.String("SET #dr = :dr"),
		ExpressionAttributeNames: map[string]string{
			"#dr": constants.DateRestoredAtt,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dr": &types.AttributeValueMemberS{
				Value: prefixItem.DateRestored,
			},
		},
	})
	// return errors
	if err != nil {
		sos.Errorf("update item date restored error: %v", err)
		return err
	}
	return nil
}

// get prefix item
func (m *_Manifest) GetPrefixItem(key interface{}) (map[string]types.AttributeValue, bool) {
	// type assert key is type item key
	keyAssert, ok := key.(registerTypes.PrefixItemKey)
	// return errors
	if !ok {
		sos.Errorf("type assert to prefix item key error")
		return nil, false
	}
	sos.Debugf("keys: [%+v]", keyAssert)
	// covert key into map[string]types.AttributeValue
	convertedKey := utils.StructToMap(keyAssert)
	// get item from table
	result, err := m.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(m.tableName),
		Key:       convertedKey,
	})
	sos.Debugf("result: [%+v]", result)
	// return errors
	if err != nil {
		sos.Errorf("get prefix item error: %v", err)
		return nil, false
	}
	// if empty, return false
	if result.Item == nil {
		return nil, false
	}
	return result.Item, true
}

// put prefix item
func (m *_Manifest) PutPrefixItem(item interface{}) error {
	// type assert to prefix item
	prefixItem, ok := item.(registerTypes.PrefixItem)
	// return errors
	if !ok {
		sos.Errorf("type assert to prefix item error")
		return nil
	}
	sos.Debugf("prefix item: [%+v]", prefixItem)
	// put item to table
	_, err := m.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(m.tableName),
		Item: map[string]types.AttributeValue{
			constants.ManifestPK: &types.AttributeValueMemberS{
				Value: prefixItem.Keys.PK,
			},
			constants.ManifestSK: &types.AttributeValueMemberS{
				Value: prefixItem.Keys.SK,
			},
			constants.DateRegisteredAtt: &types.AttributeValueMemberS{
				Value: prefixItem.DateRegistered,
			},
			constants.DateRestoredAtt: &types.AttributeValueMemberS{
				Value: prefixItem.DateRestored,
			},
			constants.DeleteMarkerIdAtt: &types.AttributeValueMemberS{
				Value: prefixItem.DeleteMarkerId,
			},
		},
	})
	// return errors
	if err != nil {
		sos.Errorf("put prefix item error: %v", err)
		return err
	}
	return nil
}

// remove prefix item
func (m *_Manifest) RemovePrefixItem(key interface{}) error {
	// type assert to prefix item
	prefixItem, ok := key.(registerTypes.PrefixItemKey)
	// return errors
	if !ok {
		sos.Errorf("type assert to prefix item error")
		return errors.New("type assert to prefix item error")
	}
	sos.Debugf("prefix item: [%+v]", prefixItem)
	// remove item from table
	_, err := m.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(m.tableName),
		Key: map[string]types.AttributeValue{
			constants.ManifestPK: &types.AttributeValueMemberS{
				Value: prefixItem.PK,
			},
			constants.ManifestSK: &types.AttributeValueMemberS{
				Value: prefixItem.SK,
			},
		},
	})
	// return errors
	if err != nil {
		sos.Errorf("remove prefix item error: %v", err)
		return err
	}
	return nil
}

// get reserved prefix list
func (m *_Manifest) GetPrefixList() map[string]types.AttributeValue {
	// get reserved prefix list from table
	result, err := m.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(m.tableName),
		Key: map[string]types.AttributeValue{
			constants.ManifestPK: &types.AttributeValueMemberS{
				Value: constants.ReservedPK,
			},
			constants.ManifestSK: &types.AttributeValueMemberS{
				Value: constants.ReservedSK,
			},
		},
	})
	sos.Debugf("result: [%+v]", result)
	// return errors
	if err != nil {
		sos.Errorf("get reserved prefix list error: %v", err)
		return nil
	}
	return result.Item
}

// add prefix to prefix list
func (m *_Manifest) AddPrefixToPrefixList(prefix string) error {
	keys := utils.StructToMap(registerTypes.PrefixItemKey{
		PK: constants.ReservedPK,
		SK: prefix,
	})
	sos.Debugf("keys : [%+v]", keys)
	newItems := []string{prefix}
	sos.Debugf("new items [%+v]", newItems)
	// Append prefix to the attribute list
	updateExpression := "SET prefixes = list_append(prefixes, :prefixes)"
	expressionAttributeValues := map[string]types.AttributeValue{
		":prefixes": &types.AttributeValueMemberL{
			Value: utils.CreateStringListAttribute(newItems),
		},
	}
	// create update item input
	updateItemInput := dynamodb.UpdateItemInput{
		TableName:                 aws.String(m.tableName),
		Key:                       keys,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		// Add a ConditionExpression to ensure the prefix doesn't already exist
		ConditionExpression: aws.String("attribute_not_exists(prefixes[:prefixes])"),
	}
	// update item in table
	_, err := m.client.UpdateItem(context.Background(), &updateItemInput)
	// return errors
	if err != nil {
		sos.Errorf("add prefix to prefix list error: %v", err)
		return err
	}
	return nil
}

// put prefix list
func (m *_Manifest) PutPrefixList(prefix string) error {
	_, err := m.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(m.tableName),
		Item: map[string]types.AttributeValue{
			constants.ManifestPK: &types.AttributeValueMemberS{
				Value: constants.ReservedPK,
			},
			constants.ManifestSK: &types.AttributeValueMemberS{
				Value: constants.ReservedSK,
			},
			constants.PrefixListAtt: &types.AttributeValueMemberL{
				Value: utils.CreateStringListAttribute([]string{prefix}),
			},
		},
	})
	// return errors
	if err != nil {
		sos.Errorf("put prefix list error: %v", err)
		return err
	}
	return nil
}

// return table name
func (m *_Manifest) GetTableName() string {
	return m.tableName
}

// return prefixes
func (m *_Manifest) GetPrefixes() []string {
	return prefixListCache
}
