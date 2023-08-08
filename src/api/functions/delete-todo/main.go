package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler, err := InitializeHandler()
	if err != nil {
		log.Printf("Failed to initialize handler. Error: %v", err)
		return
	}
	lambda.Start(handler.Handle)
}
