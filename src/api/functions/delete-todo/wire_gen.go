// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go-todos/api/utils"
	"go-todos/domain"
	"go-todos/storage"
)

// Injectors from wire.go:

func InitializeHandler() (*Handler, error) {
	config := utils.NewStorageConfig()
	todoRepository, err := storage.NewTodoRepository(config)
	if err != nil {
		return nil, err
	}
	todoService := domain.NewTodoService(todoRepository)
	handler := NewHandler(todoService)
	return handler, nil
}
