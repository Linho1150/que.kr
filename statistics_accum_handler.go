//go:build statistics_accum_handler

package main

import (
	"fmt"
	"context"
	"quekr/server/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, evts events.DynamoDBEvent) (string, error) {
	svc, err := service.NewService()

	if err != nil {
		return "cannot initiate service instance", err
	}

	for _, record := range evts.Records {
		sequence := record.Change.NewImage["sequence"].String()
		fmt.Printf("======== %s =======\n", sequence);
		err = svc.AccumlateStatisticsCounter(sequence)

		if err != nil {
			return "error occurred while processing event", err
		}
	}

	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
