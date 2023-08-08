package model

type Todo struct {
	UserEmail string `dynamodbav:"user_email" json:"-"`
	Id        string `dynamodbav:"id" json:"id"`
	Title     string `dynamodbav:"title" json:"title"`
	Completed bool   `dynamodbav:"completed" json:"completed"`
}
