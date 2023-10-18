package s3deh

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/interface/manifest"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos logger.Logger
)

type S3DeleteEventHandler interface {
	HandleEvent(ctx context.Context, request interface{}) (interface{}, error)
}

type _S3DeleteEventHandler struct {
	client   *s3.Client
	manifest manifest.Manifest
}

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	sos.Infof("init s3 delete event handler")
}

func (s3deh *_S3DeleteEventHandler) HandleEvent(ctx context.Context, request interface{}) (interface{}, error) {
	// type assert to listener request
	listenerRequest, ok := request.(types.ListenerRequest)
	if !ok {
		return nil, errors.New("request not of type listener request")
	}
	sos.Debugf("listener request : [%+v]", listenerRequest)
	return nil, nil

}

func (s3deh *_S3DeleteEventHandler) getDeleteMarker(objectMetadata types.S3ObjectMetadata) (types.DeleteMarker, error) {
	// get head object for given object
	result, err := s3deh.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &objectMetadata.BucketName,
		Key:    &objectMetadata.Key,
	})
	if err != nil {
		return types.DeleteMarker{}, err
	}
	// if object is a delete marker, return it
	if result.DeleteMarker {
		return types.DeleteMarker{
			VersionId: *result.VersionId,
		}, nil
	}
	// return empty delete marker
	return types.DeleteMarker{}, nil
}

func (s3dh *_S3DeleteEventHandler) getPrefix(objectMetadata types.S3ObjectMetadata) (types.Prefix, error) {
	return types.Prefix{
		Prefix: "",
	}, nil
}

func NewS3DeleteEventHandler(ctx context.Context, tableName string) *_S3DeleteEventHandler {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	m := manifest.NewManifest(tableName)
	return &_S3DeleteEventHandler{
		client:   s3.NewFromConfig(cfg),
		manifest: m,
	}
}
