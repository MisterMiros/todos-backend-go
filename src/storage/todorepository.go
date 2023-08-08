package storage

import (
	"context"
	"go-todos/model"
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

type TodoRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewTodoRepository(cnf *storageConfig.Config) (*TodoRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &TodoRepository{
		client:    client,
		tableName: cnf.TodoTableName,
	}, nil
}

func (repository *TodoRepository) GetUserTodos(userEmail string) ([]model.Todo, error) {
	todos := []model.Todo{}
	keyEx := expression.Key("user_email").Equal(expression.Value(userEmail))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Printf("Failed to build expression for query. Error: %v\n", err)
		return todos, err
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
		return todos, err
	}

	err = attributevalue.UnmarshalListOfMaps(output.Items, &todos)
	if err != nil {
		log.Printf("Failed to unmarshal query response. Error: %v\n", err)
	}

	return todos, err
}

func (repository *TodoRepository) GetTodo(userEmail, id string) (*model.Todo, error) {
	var todo model.Todo
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

	if (output.Item == nil) {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(output.Item, &todo)
	if err != nil {
		log.Printf("Failed to unmarshal query response. Error: %v\n", err)
		return nil, err
	}

	return &todo, nil
}

func (repository *TodoRepository) CreateTodo(userEmail, title string) (*model.Todo, error) {
	todo := model.Todo{
		UserEmail: userEmail,
		Id:        uuid.NewString(),
		Title:     title,
		Completed: false,
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

func (repository *TodoRepository) UpdateTodo(todo *model.Todo) error {
	key, err := getTodoKey(todo)
	if err != nil {
		log.Printf("Failed to get key for item '%v'. Error: %s\n", todo, err)
		return err
	}

	update := expression.Set(expression.Name("title"), expression.Value(todo.Title))
	update.Set(expression.Name("completed"), expression.Value(todo.Completed))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Failed to build update expression for item '%v'. Error: %s\n", todo, err)
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
		log.Printf("Failed to update item '%v'. Error: %s\n", todo, err)
		return err
	}
	return nil
}

func (repository *TodoRepository) DeleteTodo(userEmail, id string) error {
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

func getTodoKey(todo *model.Todo) (map[string]types.AttributeValue, error) {
	return getKey(todo.UserEmail, todo.Id)
}


func getKey(userEmail, id string) (map[string]types.AttributeValue, error) {
	userEmailAttr, err := attributevalue.Marshal(userEmail)
	if err != nil {
		return nil, err
	}
	idAttr, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}
	return map[string]types.AttributeValue{"user_email": userEmailAttr, "id": idAttr}, nil
}