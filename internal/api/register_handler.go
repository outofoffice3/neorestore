package register

import (
	"context"
	"errors"

	"github.com/outofoffice3/neorestore/internal/interface/finalize"
	"github.com/outofoffice3/neorestore/internal/interface/onboard"
	"github.com/outofoffice3/neorestore/internal/interface/validate"
	"github.com/outofoffice3/neorestore/pkg/types"
)

// handle requests for register-prefix route
func Handle(ctx context.Context, request types.RegisterRequest) (types.RegisterResponse, error) {
	v := validate.NewValidator(ctx)
	// create validate request
	validateRequest := types.ValidateRequest{
		S3BucketName: request.S3BucketName,
		Prefix:       request.Prefix,
		Region:       request.Region,
	}
	// validate prefix
	result, err := v.Validate(ctx, validateRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}
	// type assert to validate response type
	valResponse, ok := result.(types.ValidateResponse)
	// return error if wrong type
	if !ok {
		return types.RegisterResponse{}, errors.New("response not of type register response")
	}

	// onboard prefix
	var onboardRequest types.OnboardRequest
	o := onboard.NewOnboarder(ctx)
	onboardResponse, err := o.Onboard(ctx, onboardRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}

	// finalize prefix
	var finalizeRequest types.FinalizeRequest
	f := finalize.NewFinalize(ctx)
	finalizeResponse, err := f.Finalize(ctx, finalizeRequest)
	if err != nil {
		return types.RegisterResponse{}, err
	}

	// prefix successfully registered, return success
	return types.RegisterResponse{
		Body: "",
	}, nil
}
