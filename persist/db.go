package persist

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func CreateDynamoDbClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{URL: "http://localhost:8000"}, nil
				},
			),
		),
	)

	if err != nil {
		log.Fatalf("cannot load sdk config: %v", err)
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	return svc, nil
}
