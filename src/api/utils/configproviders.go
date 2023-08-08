package utils

import (
	"go-todos/storage/storageConfig"
	"os"
)

func NewStorageConfig() *storageConfig.Config {
	return &storageConfig.Config{
		TodoTableName: os.Getenv("TODO_TABLE_NAME"),
	}
}
