//go:build statistics_accum_handler

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, evts events.DynamoDBEvent) (string, error) {
	var count int

	for _, _ = range evts.Records {
		count += 1
	}

	return fmt.Sprintf("%d record(s) handled", count), nil
}

func main() {
	lambda.Start(HandleRequest)
}
