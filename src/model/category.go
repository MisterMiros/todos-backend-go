package model

type Category struct {
	UserEmail string `dynamodbav:"user_email"`
	Id        string `dynamodbav:"id"`
	Name      string `dynamodbav:"name"`
	Color     string `dynamodbav:"color"`
}