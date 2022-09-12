package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jaevor/go-nanoid"
)

type DeviceType string

const (
	DeviceTypePC     = DeviceType("pc")
	DeviceTypeMobile = DeviceType("mobile")
)

var sequenceRandomGenerator = initSequenceRandomGenerator()

func initSequenceRandomGenerator() func() string {
	gen, err := nanoid.Standard(30)

	if err != nil {
		panic(err)
	}

	return gen
}

type StatisticsRawItem struct {
	Sequence string `dynamodbav:"sequence"`

	ShortKey    string    `dynamodbav:"shortKey"`
	CreatedDate time.Time `dynamodbav:"createdDate"`

	IPAddress  string `dynamodbav:"ipAddress"`
	DeviceType DeviceType
	Referer    string
}

const STATISTICS_RAW_TABLENAME = "statistics_raw"

func (o *StatisticsRawItem) FillGeneratedSequence() {
	now := time.Now()
	o.Sequence = fmt.Sprintf("%d#%s", now.UnixMicro(), sequenceRandomGenerator())
}

func (o *Service) TouchStatistics(shortKey string, createdDate time.Time, ipAddress string, referer string, deviceType DeviceType) error {
	item := &StatisticsRawItem{
		ShortKey:    shortKey,
		CreatedDate: createdDate,
		IPAddress:   ipAddress,
		Referer:     referer,
		DeviceType:  deviceType,
	}

	item.FillGeneratedSequence()
	marshal, err := attributevalue.MarshalMap(item)

	if err != nil {
		return err
	}

	_, err = o.DbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      marshal,
		TableName: aws.String(STATISTICS_RAW_TABLENAME),
	})

	if err != nil {
		return err
	}

	return nil
}
