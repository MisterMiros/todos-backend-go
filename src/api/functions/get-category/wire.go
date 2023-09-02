//go:build wireinject
// +build wireinject

package main

import (
	"go-todos/api/utils"
	"go-todos/domain"
	"go-todos/storage"

	"github.com/google/wire"
)

func InitializeHandler() (*Handler, error) {
	wire.Build(utils.NewStorageConfig, storage.NewTodoRepository, storage.NewCategoryRepository, domain.NewCategoryService, NewHandler)
	return &Handler{}, nil
}
