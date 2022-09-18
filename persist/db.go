package persist

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func checkUseLocalDynamoDbFlag() bool {
	useLocalDynamoDbFlagStr := os.Getenv("USELOCALDYNAMODB")

	if useLocalDynamoDbFlagStr != "" {
		return true
	}

	return false
}

func CreateDynamoDbClient() (*dynamodb.Client, error) {
	var cfgOptFns []func(*config.LoadOptions) error

	if checkUseLocalDynamoDbFlag() {
		cfgOptFns = append(
			cfgOptFns,
			config.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(
					func(service, region string, options ...interface{}) (aws.Endpoint, error) {
						return aws.Endpoint{URL: "http://localhost:8000"}, nil
					},
				),
			),
		)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		cfgOptFns...,
	)

	if err != nil {
		log.Fatalf("cannot load sdk config: %v", err)
		return nil, err
	}

	svc := dynamodb.NewFromConfig(cfg)
	return svc, nil
}
