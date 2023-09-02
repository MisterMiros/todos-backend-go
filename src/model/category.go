package model

type Category struct {
	UserEmail string `dynamodbav:"user_email" json:"-"`
	Id        string `dynamodbav:"id" json:"id"`
	Name      string `dynamodbav:"name" json:"name"`
	Color     string `dynamodbav:"color" json:"color"`
}