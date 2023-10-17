package onboard

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/outofoffice3/neorestore/pkg/constants"
	"github.com/outofoffice3/neorestore/pkg/types"
)

type Onboarder interface {
	// complete any activities to properly onboard prefix to be monitored
	Onboard(ctx context.Context, request interface{}) (interface{}, error)
}

type _Onboarder struct {
	client *s3.Client
}

func (o *_Onboarder) Onboard(ctx context.Context, request interface{}) (interface{}, error) {
	// type assert request to onboard request
	onboardRequest, ok := request.(types.OnboardRequest)
	if !ok {
		return types.OnboardResponse{}, nil
	}
	// create s3 delete notification configuration for bucket
	err := o.createDeleteNotificationConfig(onboardRequest.S3BucketName)
	if err != nil {
		return types.OnboardResponse{}, err
	}
	// return success response
	return types.OnboardResponse{
		Body: "S3 delete configuration successfully created",
	}, nil
}

// create s3 delete notification configuation
func (o *_Onboarder) createDeleteNotificationConfig(bucketName string) error {
	functionArn := constants.ListenerFunctionArn
	// create s3 delete configuration
	config := s3Types.NotificationConfiguration{
		LambdaFunctionConfigurations: []s3Types.LambdaFunctionConfiguration{
			{
				LambdaFunctionArn: aws.String(functionArn),
				Events:            []s3Types.Event{s3Types.EventS3ObjectRemoved},
			},
		},
	}
	// create input
	input := &s3.PutBucketNotificationConfigurationInput{
		Bucket:                    aws.String(bucketName),
		NotificationConfiguration: &config,
		SkipDestinationValidation: true,
	}
	// make call to s3
	_, err := o.client.PutBucketNotificationConfiguration(context.Background(), input)
	if err != nil {
		return err
	}
	return nil
}

// create new onboarder
func NewOnboarder(ctx context.Context) Onboarder {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	return &_Onboarder{
		client: s3.NewFromConfig(cfg),
	}
}
