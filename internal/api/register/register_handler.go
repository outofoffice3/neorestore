package register

import (
	"context"
	"errors"

	"github.com/outofoffice3/common/logger"
	"github.com/outofoffice3/neorestore/internal/interface/finalize"
	"github.com/outofoffice3/neorestore/internal/interface/onboard"
	"github.com/outofoffice3/neorestore/internal/interface/validate"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos       logger.Logger
	validator validate.Validator
	onboarder onboard.Onboarder
	finalizer finalize.Finalizer
)

func init() {
	sos = logger.NewConsoleLogger(logger.LogLevelInfo)
	sos.Infof("init register handler")
	// create new validator
	validator = validate.NewValidator(context.Background())
	// create new onboarder
	onboarder = onboard.NewOnboarder(context.Background())
	// create new finalizer
	finalizer = finalize.NewFinalizer(context.Background())
	sos.Infof("init register handler completed")
}

// handle requests for register-prefix route
func Handle(ctx context.Context, request types.RegisterRequest) (types.RegisterResponse, error) {
	// create validate request
	validateRequest := types.ValidateRequest{
		S3BucketName: request.S3BucketName,
		Prefix:       request.Prefix,
		Region:       request.Region,
	}
	// validate prefix
	result, err := validator.Validate(ctx, validateRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}
	// type assert to validate response type
	valResponse, ok := result.(types.ValidateResponse)
	// return error if wrong type
	if !ok {
		return types.RegisterResponse{}, errors.New("response not of type register response")
	}
	sos.Debugf("validator response : %+v", valResponse.Body)

	// onboard prefix
	onboardRequest := types.OnboardRequest{
		S3BucketName: request.S3BucketName,
		Prefix:       request.Prefix,
		Region:       request.Region,
	}
	result, err = onboarder.Onboard(ctx, onboardRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}
	// type assert to onboard response type
	onboardResponse, ok := result.(types.OnboardResponse)
	if !ok {
		return types.RegisterResponse{}, errors.New("response not of type onboard response")
	}
	sos.Debugf("onboard response : %+v", onboardResponse.Body)

	// finalize prefix
	finalizeRequest := types.FinalizeRequest{
		S3BucketName: request.S3BucketName,
		Prefix:       request.Prefix,
		Region:       request.Region,
	}
	f := finalize.NewFinalizer(ctx)
	result, err = f.Finalize(ctx, finalizeRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}
	// type assert to finalize response type
	finalizeResponse, ok := result.(types.FinalizeResponse)
	if !ok {
		return types.RegisterResponse{}, errors.New("response not of type finalize response")
	}
	sos.Debugf("finalize response : %+v", finalizeResponse.Body)

	// prefix successfully registered, return success
	return types.RegisterResponse{
		Body: "",
	}, nil
}
