package validate

import (
	"context"
	"errors"

	"github.com/outofoffice3/neorestore/internal/interface/manifest"
	"github.com/outofoffice3/neorestore/pkg/types"
	"github.com/outofoffice3/neorestore/pkg/utils"
)

type Validator interface {
	Validate(ctx context.Context, request interface{}) (interface{}, error)
}

type _Validator struct {
	manifest manifest.Manifest
}

// handle requests for register-prefix route
func (v *_Validator) Validate(ctx context.Context, request interface{}) (interface{}, error) {
	// type assert request to make sure it is RegisterRequest type
	requestAssert, ok := request.(types.ValidateRequest)
	if !ok {
		return types.ValidateResponse{}, errors.New("request not type register request")
	}
	// validate prefix
	validateResult, err := v.validatePrefix(ctx, requestAssert)
	// if prefix is invalid, return error
	if !validateResult {
		return types.RegisterResponse{}, err
	}
	// prefix is valid, return success
	prefix := requestAssert.S3BucketName + "/" + requestAssert.Prefix
	return types.RegisterResponse{
		Body: "Prefix " + prefix + " is successfully registered",
	}, nil
}

// validate if a prefix is valid
func (v *_Validator) validatePrefix(ctx context.Context, request types.ValidateRequest) (bool, error) {
	// get prefix from request
	prefix := request.S3BucketName + "/" + request.Prefix
	// check if prefix exists
	result := v.manifest.GetPrefixList()
	// if prefix exists, return error
	if utils.ContainsPrefix(result, prefix) {
		return false, errors.New("prefix already exists")
	}
	// prefix does not exist, return true for valid prefix
	return true, nil
}

// create a new Validator
func NewValidator(ctx context.Context) Validator {
	manifest := manifest.NewManifest()
	return &_Validator{
		manifest: manifest,
	}
}
