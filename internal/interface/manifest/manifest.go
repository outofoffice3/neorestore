package manifest

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/outofoffice3/common/logger"
)

type Manifest interface {
	// get a prefix item from table
	GetPrefixItem(key interface{}) (map[string]types.AttributeValue, bool)
	// put a prefix item to table
	PutPrefixItem(item interface{}) error
	// remove prefix item from table
	RemovePrefixItem(key interface{}) error
	// get reserved prefix item from table
	GetPrefixList() map[string]types.AttributeValue
	// put reserved prefix item to table
	AddPrefixToPrefixList(prefix string) error
}

var (
	sos logger.Logger
)

type _Manifest struct {
	client    *dynamodb.Client
	tableName string
}

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelInfo)
	sos.Infof("manifest init completed")
}

// create new Manifest
func NewManifest(tableName string) _Manifest {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)

	}
	return _Manifest{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
}

// get prefix item
func (m _Manifest) GetPrefixItem(key interface{}) (map[string]types.AttributeValue, bool) {
	return nil, false
}

// put prefix item
func (m _Manifest) PutPrefixItem(item interface{}) error {
	return nil
}

// remove prefix item
func (m _Manifest) RemovePrefixItem(key interface{}) error {
	return nil
}

// get reserved prefix list
func (m _Manifest) GetPrefixList() map[string]types.AttributeValue {
	return nil
}

// add prefix to prefix list
func (m _Manifest) AddPrefixToPrefixList(prefix string) error {
	return nil
}
