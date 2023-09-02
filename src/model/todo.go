package model

type Todo struct {
	UserEmail  string `dynamodbav:"user_email"`
	Id         string `dynamodbav:"id"`
	Title      string `dynamodbav:"title"`
	Completed  bool   `dynamodbav:"completed"`
	CategoryId string `dynamodbav: "category_id"`
}
