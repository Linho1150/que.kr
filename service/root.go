package service

import (
	"quekr/server/persist"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Service struct {
	DbClient      *dynamodb.Client
	LocalTimezone *time.Location
}

func NewService() (*Service, error) {
	client, err := persist.CreateDynamoDbClient()

	if err != nil {
		return nil, err
	}

	location, err := time.LoadLocation("Asia/Seoul")

	if err != nil {
		return nil, err
	}

	return &Service{
		DbClient:      client,
		LocalTimezone: location,
	}, nil
}

func (o *Service) NowLocalTime() time.Time {
	time := time.Now().In(o.LocalTimezone)
	return time
}
