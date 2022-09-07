package service

import (
	"quekr/server/persist"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Service struct {
	DbClient *dynamodb.Client
}

func NewService() (*Service, error) {
	client, err := persist.CreateDynamoDbClient()

	if err != nil {
		return nil, err
	}

	return &Service{
		DbClient: client,
	}, nil
}
