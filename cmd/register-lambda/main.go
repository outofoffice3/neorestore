package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	register "github.com/outofoffice3/neorestore/internal/api/register"
	"github.com/outofoffice3/neorestore/pkg/types"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// switch based on request path
	path := request.Path
	switch path {
	case "/register-prefix":
		// convert request body to RegisterPrefixRequest
		var registerRequest types.RegisterRequest
		err := json.Unmarshal([]byte(request.Body), &registerRequest)
		// return errors if marshal fails
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}
		// register prefifx
		result, err := register.Handle(ctx, registerRequest)
		// return errors if register fails
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500,
				Body: err.Error(),
			}, nil
		}
		// return success
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       result.Body,
		}, nil
	default:
		{

		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       "catch all",
	}, nil
}

func main() {
	lambda.Start(handler)
}
