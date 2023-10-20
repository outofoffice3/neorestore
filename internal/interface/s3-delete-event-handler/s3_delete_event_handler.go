package s3deh

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/interface/manifest"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos         logger.Logger
	cachedS3deh S3DeleteEventHandler
)

type S3DeleteEventHandler interface {
	// handle s3 delete events
	HandleEvent(ctx context.Context, request interface{}) (interface{}, error)
}

type _S3DeleteEventHandler struct {
	client   *s3.Client        // s3 client
	manifest manifest.Manifest // manifest client
}

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	// create new instance of s3deh
	cachedS3deh = NewS3DeleteEventHandler(context.Background())
	sos.Infof("init s3 delete event handler")
}

// set logger level
func SetLogLevel(level logger.LogLevel) {
	sos.SetLogLevel(level)
}

func (s3deh *_S3DeleteEventHandler) HandleEvent(ctx context.Context, event interface{}) (interface{}, error) {
	// type assert to listener event
	le, ok := event.(types.ListenerEvent)
	if !ok {
		return nil, errors.New("event not of type listener event")
	}
	sos.Debugf("listener event : [%+v]", le)
	// get the delete marker for the deleted object
	deleteMarker, err := s3deh.getDeleteMarker(le.S3ObjectMetadata)
	if err != nil {
		return nil, err
	}
	sos.Debugf("delete marker : [%+v]", deleteMarker)
	// get the matching prefix for the object
	prefix, err := s3deh.getPrefix(le.S3ObjectMetadata)
	if err != nil {
		return nil, err
	}
	sos.Debugf("prefix : [%+v]", prefix)
	// if the event object and the delete are the same
	// handle it as a restore event
	var prefixItem types.PrefixItem
	if isRestoreEvent(le.S3ObjectMetadata, deleteMarker) {
		err := s3deh.handleItemRestoreEvent(prefixItem)
		if err != nil {
			return nil, err
		}
		// return success
		return types.ListenerResponse{
			Body: "success",
		}, nil
	}
	// if the event object and the delete are different
	// handle it as a new delete marker event
	err = s3deh.handleNewDeleteMarkerEvent(prefixItem)
	if err != nil {
		return nil, err
	}
	// return success
	return types.ListenerResponse{
		Body: "success",
	}, nil
}

// handle item restore events
// update date restored attribute for given item
func (s3deh *_S3DeleteEventHandler) handleItemRestoreEvent(item types.PrefixItem) error {
	// update item in manifest
	err := s3deh.manifest.UpdateItemDateRestored(item)
	if err != nil {
		return err
	}
	return nil
}

// handle new delete marker event
func (s3deh *_S3DeleteEventHandler) handleNewDeleteMarkerEvent(item types.PrefixItem) error {
	// put item in manifest
	err := s3deh.manifest.PutPrefixItem(item)
	if err != nil {
		return err
	}
	return nil
}

// function to see if the event is a restore event
func isRestoreEvent(objectMetadata types.S3ObjectMetadata, deleteMarker types.DeleteMarker) bool {
	return objectMetadata.VersionId == deleteMarker.VersionId
}

// get delete marker for a deleted s3 object
func (s3deh *_S3DeleteEventHandler) getDeleteMarker(objectMetadata types.S3ObjectMetadata) (types.DeleteMarker, error) {
	// get head object for given object
	result, err := s3deh.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &objectMetadata.BucketName,
		Key:    &objectMetadata.Key,
	})
	sos.Debugf("head object result : [%+v]", result)
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
	return types.DeleteMarker{
		VersionId: "",
	}, nil
}

// get matching prefix from manifest for given object
func (s3dh *_S3DeleteEventHandler) getPrefix(objectMetadata types.S3ObjectMetadata) (types.Prefix, error) {
	//  get prefix list from manifest
	var itemMap map[string]interface{}
	result := s3dh.manifest.GetPrefixList()
	err := attributevalue.UnmarshalMap(result, &itemMap)
	if err != nil {
		return types.Prefix{
			Prefix: "",
		}, err
	}
	// read prefixes attribute
	prefixes, ok := itemMap["prefixes"].([]string)
	if !ok {
		return types.Prefix{
			Prefix: "",
		}, errors.New("prefixes attribute not of type []string")
	}
	// loop through prefixes and find the one that is contained in the bucket name
	// + object key of the event object
	fullObjKey := objectMetadata.BucketName + objectMetadata.Key
	for _, prefix := range prefixes {
		// check if prefix is contained inside of full object key
		sos.Debugf("full object key [%v] & prefix [%v]", fullObjKey, prefix)
		if strings.Contains(fullObjKey, prefix) {
			return types.Prefix{
				Prefix: prefix,
			}, nil
		}
	}
	return types.Prefix{
		Prefix: "",
	}, nil
}

// create new s3 delete event handler
func NewS3DeleteEventHandler(ctx context.Context) S3DeleteEventHandler {
	// return cached version if availabe
	if cachedS3deh != nil {
		return cachedS3deh
	}
	// if nil, create new instance
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	m := manifest.NewManifest()
	// update cached version
	cachedS3deh = &_S3DeleteEventHandler{
		client:   s3.NewFromConfig(cfg),
		manifest: m,
	}
	// return new instance
	return cachedS3deh
}
