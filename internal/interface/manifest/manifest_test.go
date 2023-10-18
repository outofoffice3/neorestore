package manifest

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/joho/godotenv"
	"github.com/outofoffice3/neorestore/pkg/types"
	"github.com/stretchr/testify/assert"
)

// test cases for manifest interface

func TestNewManifest(t *testing.T) {
	// init manifest package
	Init()
	a := assert.New(t)
	// get table name from .env
	filePath := "../../../.env"
	err := godotenv.Load(filePath)
	// assert error is nil
	a.Nil(err, "error is not nil")
	tableName := os.Getenv("manifestTableName")
	// create manifest
	m := NewManifest(tableName)
	// assert manifest is not nil
	a.NotNil(m, "manifest is nil")
	// assert manifest table name is equal to table name
	a.Equal(tableName, m.GetTableName(), "manifest table name is not equal to table name")
}

func TestPrefixList(t *testing.T) {
	// init manifest package
	Init()
	a := assert.New(t)
	// get table name from .env
	filePath := "../../../.env"
	err := godotenv.Load(filePath)
	// assert error is nil
	a.Nil(err, "error is not nil")
	tableName := os.Getenv("manifestTableName")
	// create manifest
	m := NewManifest(tableName)
	// put prefix list
	err = m.PutPrefixList("test-prefix")
	// assert error is nil
	a.Nil(err, "error is not nil")
	// get prefix list
	item := m.GetPrefixList()
	// assert prefix list is not nil
	a.NotNil(item, "prefix list is nil")
	// assert prefix list is equal to "test-prefix"
	prefixes := item["prefixes"]
	// assert prefixes attribute is not nil
	a.NotNil(prefixes, "prefixes attribute is nil")
	// unmarshall to prefix list type
	var prefixListItem types.PrefixListItem
	err = attributevalue.UnmarshalMap(item, &prefixListItem)
	// assert error is nil
	a.Nil(err, "error is not nil")
	// assert prefix list is equal to "test-prefix"
	a.Equal("test-prefix", prefixListItem.Prefixes[0], "prefix list is not equal to test-prefix")

}

func TestPrefixItem(t *testing.T) {
	// init manifest package
	Init()
	a := assert.New(t)
	// get table name from .env
	filePath := "../../../.env"
	err := godotenv.Load(filePath)
	// assert error is nil
	a.Nil(err, "error is not nil")
	tableName := os.Getenv("manifestTableName")
	// create manifest
	m := NewManifest(tableName)
	// put prefix item
	item := types.PrefixItem{
		Keys: types.PrefixItemKey{
			PK: "primarykey",
			SK: "sortkey",
		},
		DateRegistered: "testregistered",
		DateRestored:   "testrestored",
		DeleteMarkerId: "testdeletemarkerid",
	}
	err = m.PutPrefixItem(item)
	// assert error is nil
	a.Nil(err, "error is not nil")
	// get prefix item
	prefixItemKey := types.PrefixItemKey{
		PK: "primarykey",
		SK: "sortkey",
	}
	prefixItem, ok := m.GetPrefixItem(prefixItemKey)
	var pi types.PrefixItem
	err = attributevalue.UnmarshalMap(prefixItem, &pi)
	// assert error is nil
	a.NoError(err, "error is not nil")
	// assert prefix item is found
	a.True(ok, "prefix item should be found")
	// assert prefix item is not nil
	a.NotNil(prefixItem, "prefix item is nil")
	// assert prefix item date registered is equal to "testregistered"
	a.Equal("testregistered", pi.DateRegistered, "prefix item date registered is not equal to testregistered")
	// assert prefix item date restored is equal to "testrestored"
	a.Equal("testrestored", pi.DateRestored, "prefix item date restored is not equal to testrestored")
	// assert prefix item delete marker id is equal to "testdeletemarkerid"
	a.Equal("testdeletemarkerid", pi.DeleteMarkerId, "prefix item delete marker id is not equal to testdeletemarkerid")
	// delete prefix item
	err = m.RemovePrefixItem(prefixItemKey)
	// assert error is nil
	a.Nil(err, "error is not nil")
	// get prefix item
	prefixItem, ok = m.GetPrefixItem(prefixItemKey)
	// assert prefix item is not found
	a.False(ok, "prefix item should not be found")
	// assert prefix item is nil
	a.Nil(prefixItem, "prefix item is not nil")
}
