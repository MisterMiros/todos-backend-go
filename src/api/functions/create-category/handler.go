package main

import (
	"encoding/json"
	"go-todos/api/apimodels"
	"go-todos/api/utils"
	"go-todos/api/utils/responses"
	"go-todos/domain"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	categoryService *domain.CategoryService
}

func NewHandler(categoryService *domain.CategoryService) *Handler {
	return &Handler{
		categoryService: categoryService,
	}
}

func (handler *Handler) Handle(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	email, err := utils.GetEmail(event)
	if err != nil {
		log.Printf("Failed to parse authorizer context. Error: %v", err)
		return responses.BadRequest("Failed to parse authorizer context")
	}

	var body apimodels.CreateCategoryRequest
	err = json.Unmarshal([]byte(event.Body), &body)
	if err != nil {
		log.Printf("Failed to parse request body. Error: %v", err)
		return responses.InternalError("Failed to parse request body")
	}
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" {
		return responses.BadRequest("Can't create item with an empty title")
	}

	category, serviceErr := handler.categoryService.CreateCategory(email, body.Name, body.Color)
	if serviceErr != nil {
		log.Printf("Failed to create category. Error: %v", serviceErr)
		return utils.ResponseFromServiceError(serviceErr)
	}
	return responses.Ok(apimodels.NewCategory(category))
}
