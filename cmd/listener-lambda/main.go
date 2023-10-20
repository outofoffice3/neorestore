package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/api/listener"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos logger.Logger
)

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	sos.Infof("init s3 event handler")
}

func handler(ctx context.Context, event events.S3Event) error {
	// delete event comes from s3 bucket
	for _, record := range event.Records {
		// save s3 object metadata into custom struct
		objMetadata := types.S3ObjectMetadata{
			BucketName: record.S3.Bucket.Name,
			Key:        record.S3.Object.Key,
			VersionId:  record.S3.Object.VersionID,
			Region:     record.AWSRegion,
		}
		// save metadata into custom struct
		eventMetadata := types.EventMetadata{
			Name:    record.EventName,
			Time:    record.EventName,
			Version: record.EventVersion,
			Region:  record.AWSRegion,
		}
		le := types.ListenerEvent{
			S3ObjectMetadata: objMetadata,
			EventMetadata:    eventMetadata,
		}
		result, err := listener.Handle(ctx, le)
		// return errors
		if err != nil {
			return err
		}
		sos.Debugf("Listener response [%+v]", result)

	}
	return nil
}

func main() {
	lambda.Start(handler)
}
