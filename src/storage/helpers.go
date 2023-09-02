package storage

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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