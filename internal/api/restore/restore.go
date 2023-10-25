package restore

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/labstack/echo/v4"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/interface/manifest"
	"github.com/outofoffice3/neorestore/pkg/constants"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos logger.Logger
	mt  manifest.Manifest
)

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	sos.Infof("init restore api successfull")
	mt = manifest.NewManifest()
}

// create handler for restore requests
func RestoreHandler(c echo.Context) error {
	var (
		requestData types.RestoreRequest
		items       []map[string]ddbTypes.AttributeValue
	)
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid json data")
	}
	sos.Debugf("request data : [%+v]", requestData)
	prefix := requestData.Prefix
	// query items from dynamoDB for a the given prefix
	keyConditionExpression := "#pk = :pk"
	expressionAttributeNames := map[string]string{
		"#pk": constants.ManifestPK,
	}
	expressionAttributeValues := map[string]ddbTypes.AttributeValue{
		":pk": &ddbTypes.AttributeValueMemberS{
			Value: prefix,
		},
	}
	queryInput := dynamodb.QueryInput{
		TableName:                 aws.String(mt.GetTableName()),
		KeyConditionExpression:    &keyConditionExpression,
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	// get query results
	result := mt.Query(&queryInput)
	// loop through query results
	for result.HasMorePages() {
		// save query output
		output, err := result.NextPage(context.Background())
		// return errors
		if err != nil {
			sos.Errorf("failed to get next page of query results: %s", err)
			return c.JSON(http.StatusInternalServerError, "failed to get next page of query results")
		}
		// save items from query output
		items = output.Items
		// loop through items
		for _, item := range items {
			sos.Debugf("item : [%+v]", item)
		}
	}
	return c.JSON(http.StatusOK, requestData)
}
