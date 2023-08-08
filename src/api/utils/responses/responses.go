package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type errorBody struct {
	Message string `json:"message"`
}

func Ok(body interface{}) (*events.APIGatewayProxyResponse, error) {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal response body. Error: %v\n", err)
		return nil, err
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(bodyJson),
	}, nil
}

func NoContent() (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func BadRequest(message string) (*events.APIGatewayProxyResponse, error) {
	return statusCodeError(&message, http.StatusBadRequest)
}

func InternalError(message string) (*events.APIGatewayProxyResponse, error) {
	return statusCodeError(&message, http.StatusInternalServerError)
}

func NotFound() (*events.APIGatewayProxyResponse, error) {
	return statusCodeError(nil, http.StatusNotFound)
}

func statusCodeError(message *string, statusCode int) (*events.APIGatewayProxyResponse, error) {
	var bodyJson string
	if message != nil {
		body := errorBody{
			Message: *message,
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			log.Printf("Failed to marshal error body. Error: %v\n", err)
			return nil, err
		}
		bodyJson = string(bodyBytes)
	}
	
	log.Printf("Creating response...\n")
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       bodyJson,
	}, nil
}
