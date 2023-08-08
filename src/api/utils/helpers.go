package utils

import (
	"fmt"
	"go-todos/api/utils/responses"
	"go-todos/domain"

	"github.com/aws/aws-lambda-go/events"
)

func GetEmail(request events.APIGatewayProxyRequest) (string, error) {
	authorizer, ok := request.RequestContext.Authorizer["jwt"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not assert Authorizer['jwt'] as map[string]interface{}")
	}

	claims, ok := authorizer["claims"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not assert claims as map[string]interface{}")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("could not assert email as string")
	}

	return email, nil
}

func GetId(request events.APIGatewayProxyRequest) (string, error) {
	id, ok := request.PathParameters["id"]
	if !ok {
		return "", fmt.Errorf("'id' parameter is missing from path")
	}
	return id, nil
}

func ResponseFromServiceError(err *domain.ServiceError) (*events.APIGatewayProxyResponse, error) {
	switch err.Kind {
	case domain.InternalError:
		return responses.InternalError(err.Message.Error())
	case domain.NotFound:
		return responses.NotFound()
	}
	panic("Unknown error kind")
}
