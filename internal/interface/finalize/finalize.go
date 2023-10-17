package finalize

import (
	"context"
	"errors"

	"github.com/outofoffice3/neorestore/internal/interface/manifest"
	"github.com/outofoffice3/neorestore/pkg/types"
)

type Finalizer interface {
	// finalize the registration of the prefix
	Finalize(ctx context.Context, request interface{}) (interface{}, error)
}

type _Finalize struct {
	manifest manifest.Manifest
}

func (f *_Finalize) Finalize(ctx context.Context, request interface{}) (interface{}, error) {
	// type assert to finalize request type
	finalizeRequest, ok := request.(types.FinalizeRequest)
	if !ok {
		return nil, errors.New("request not of type FinalizeRequest")
	}
	// add prefix to prefix list
	prefix := finalizeRequest.S3BucketName + "/" + finalizeRequest.Prefix
	err := f.manifest.AddPrefixToPrefixList(prefix)
	if err != nil {
		return nil, err
	}
	return types.FinalizeResponse{
		Body: "Prefix added to prefix list",
	}, nil
}

func NewFinalizer(ctx context.Context) Finalizer {
	m := manifest.NewManifest("")
	return &_Finalize{
		manifest: m,
	}
}
