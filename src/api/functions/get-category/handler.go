package main

import (
	"go-todos/api/apimodels"
	"go-todos/api/utils"
	"go-todos/api/utils/responses"
	"go-todos/domain"
	"log"

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

	id, err := utils.GetId(event)
	if err != nil {
		log.Printf("Failed to get category 'id'. Error: %v", err)
		return responses.BadRequest("Failed to get category 'id'")
	}

	category, serviceErr := handler.categoryService.GetCategory(email, id)
	if serviceErr != nil {
		log.Printf("Failed to get category. Error: %v", serviceErr)
		return utils.ResponseFromServiceError(serviceErr)
	}
	return responses.Ok(apimodels.NewCategory(category))
}
