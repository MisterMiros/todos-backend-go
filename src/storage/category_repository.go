package storage

import (
	"context"
	"go-todos/model"
	"go-todos/storage/interfaces"
	"go-todos/storage/storageConfig"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type CategoryRepository struct {
	client    interfaces.DynamoClient
	tableName string
}

func NewCategoryRepository(cnf storageConfig.Config) (*CategoryRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &CategoryRepository{
		client:    client,
		tableName: cnf.CategoryTableName,
	}, nil
}

func (repository *CategoryRepository) GetUserCategories(userEmail string) ([]model.Category, error) {
	categories := []model.Category{}
	keyEx := expression.Key("user_email").Equal(expression.Value(userEmail))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Printf("Failed to build expression for query. Error: %v\n", err)
		return categories, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(repository.tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	output, err := repository.client.Query(context.TODO(), input)

	if err != nil {
		log.Printf("Failed to query items for user '%v'. Error: %v\n", userEmail, err)
		return categories, err
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, &categories)
	if err != nil {
		log.Printf("Failed to unmarshal query response. Error: %v\n", err)
	}

	return categories, err
}

func (repository *CategoryRepository) GetCategory(userEmail, id string) (*model.Category, error) {
	var todo model.Category
	key, err := getKey(userEmail, id)
	if err != nil {
		log.Printf("Failed to get key for item '%s:%s'. Error: %s\n", userEmail, id, err)
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(repository.tableName),
		Key:       key,
	}

	output, err := repository.client.GetItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to get item '%s:%s'. Error: %s\n", userEmail, id, err)
		return nil, err
	}

	if output.Item == nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(output.Item, &todo)
	if err != nil {
		log.Printf("Failed to unmarshal query response. Error: %v\n", err)
		return nil, err
	}

	return &todo, nil
}

func (repository *CategoryRepository) CreateCategory(userEmail, name string, color string) (*model.Category, error) {
	todo := model.Category{
		UserEmail: userEmail,
		Id:        uuid.NewString(),
		Name:      name,
		Color:     color,
	}

	item, err := attributevalue.MarshalMap(todo)
	if err != nil {
		log.Printf("Failed to marshal item '%v'. Error: %s\n", todo, err)
		return nil, err
	}

	input := dynamodb.PutItemInput{
		TableName: aws.String(repository.tableName),
		Item:      item,
	}
	_, err = repository.client.PutItem(context.TODO(), &input)
	if err != nil {
		log.Printf("Failed to put item '%v'. Error: %s\n", todo, err)
		return nil, err
	}
	return &todo, nil
}

func (repository *CategoryRepository) UpdateCategory(category *model.Category) error {
	key, err := getCategoryKey(category)
	if err != nil {
		log.Printf("Failed to get key for item '%v'. Error: %s\n", category, err)
		return err
	}

	update := expression.Set(expression.Name("name"), expression.Value(category.Name))
	update.Set(expression.Name("color"), expression.Value(category.Color))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Failed to build update expression for item '%v'. Error: %s\n", category, err)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(repository.tableName),
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	}
	_, err = repository.client.UpdateItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to update item '%v'. Error: %s\n", category, err)
		return err
	}
	return nil
}

func (repository *CategoryRepository) DeleteCategory(userEmail, id string) error {
	key, err := getKey(userEmail, id)
	if err != nil {
		log.Printf("Failed to get key for item '%s:%s'. Error: %s\n", userEmail, id, err)
		return err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(repository.tableName),
		Key:       key,
	}

	_, err = repository.client.DeleteItem(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to delete item '%s:%s'. Error: %s\n", userEmail, id, err)
		return err
	}
	return nil
}

func getCategoryKey(category *model.Category) (map[string]types.AttributeValue, error) {
	return getKey(category.UserEmail, category.Id)
}
