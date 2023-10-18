package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/api/listener"
	s3deleteeventhandler "github.com/outofoffice3/neorestore/internal/interface/s3-delete-event-handler"
	"github.com/outofoffice3/neorestore/pkg/constants"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos  logger.Logger
	s3dh s3deleteeventhandler.S3DeleteEventHandler
)

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	s3dh = s3deleteeventhandler.NewS3DeleteEventHandler(context.Background(), constants.ManifestTableName)
	sos.Infof("init s3 event handler i")
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
		listenerRequest := types.ListenerRequest{
			S3ObjectMetadata: objMetadata,
			EventMetadata:    eventMetadata,
		}
		result, err := listener.Handle(ctx, listenerRequest)
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
