package main

import (
	"encoding/json"
	"go-todos/api/apimodels"
	"go-todos/api/utils"
	"go-todos/api/utils/responses"
	"go-todos/domain"
	"go-todos/model"
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

	id, err := utils.GetId(event)
	if err != nil {
		log.Printf("Failed to get category 'id'. Error: %v", err)
		return responses.BadRequest("Failed to get category 'id'")
	}

	var body apimodels.UpdateCategoryRequest
	err = json.Unmarshal([]byte(event.Body), &body)
	if err != nil {
		log.Printf("Failed to parse request body. Error: %v", err)
		return responses.InternalError("Failed to parse request body")
	}
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" {
		return responses.BadRequest("Can't set empty name to category")
	}

	category := model.Category {
		UserEmail: email,
		Id: id,
		Name: body.Name,
		Color: body.Color,
	}

	serviceErr := handler.categoryService.UpdateCategory(&category)
	if serviceErr != nil {
		log.Printf("Failed to update category. Error: %v", serviceErr)
		return utils.ResponseFromServiceError(serviceErr)
	}
	return responses.NoContent()
}
