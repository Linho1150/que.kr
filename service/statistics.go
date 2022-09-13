package service

import (
	"context"
	"errors"
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

type StatisticsAccumDeviceTypeItem struct {
	ShortKey   string     `dynamodbav:"shortKey"`
	DeviceType DeviceType `dynamodbav:"devicetype"`
	Counter    int        `dynamodbav:"counter"`
}

type StatisticsAccumRefererItem struct {
	ShortKey string `dynamodbav:"shortKey"`
	Referer  string `dynamodbav:"referer"`
	Counter  int    `dynamodbav:"counter"`
}

type StatisticsAccumTimePerDate struct {
	ShortKey string    `dynamodbav:"shortKey"`
	DateTime time.Time `dynamodbav:"datetime"`
	Counter  int       `dynamodbav:"counter"`
}

type StatisticsAccumTimePerMinute struct {
	ShortKey string    `dynamodbav:"shortKey"`
	DateTime time.Time `dynamodbav:"datetime"`
	Counter  int       `dynamodbav:"counter"`
}

var STATISTICS_RAW_TABLENAME = "statistics_raw"

func (o *Service) AccumlateStatisticsCounter(sequence string) error {
	key := map[string]string{
		"sequence": sequence,
	}

	marshalKey, err := attributevalue.MarshalMap(key)

	if err != nil {
		return err
	}

	result, err := o.DbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &STATISTICS_RAW_TABLENAME,
		Key:       marshalKey,
	})

	if err != nil {
		return err
	}

	if result.Item == nil {
		return errors.New("specified sequence item does not exist")
	}

	item := &StatisticsRawItem{}

	if err = attributevalue.UnmarshalMap(result.Item, item); err != nil {
		return err
	}

	o.IncrementStatisticsCounter("time_per_date", item.ShortKey, "datetime", RoundDateTime(&item.CreatedDate, 60*60*24))
	o.IncrementStatisticsCounter("time_per_minute", item.ShortKey, "datetime", RoundDateTime(&item.CreatedDate, 60))
	o.IncrementStatisticsCounter("referer", item.ShortKey, "referer", item.Referer)
	o.IncrementStatisticsCounter("devicetype", item.ShortKey, "devicetype", item.DeviceType)
	return nil
}

func RoundDateTime(target *time.Time, seconds int) *time.Time {
	secs := (int64(target.UnixMilli()) / int64(seconds)) * int64(seconds)
	ret := time.UnixMilli(secs)
	return &ret
}

func (o *Service) IncrementStatisticsCounter(tableName string, shortKey string, secondaryKeyName string, secondaryKeyValue interface{}) error {
	fullTableName := aws.String("statistics_accum_" + tableName)

	key := map[string]interface{}{
		"shortKey":       shortKey,
		secondaryKeyName: secondaryKeyValue,
	}

	marshalKey, err := attributevalue.MarshalMap(key)

	if err != nil {
		return err
	}

	result, err := o.DbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: fullTableName,
		Key:       marshalKey,
	})

	if err != nil {
		return err
	}

	newItem := make(map[string]interface{})

	if result.Item != nil {
		if err = attributevalue.UnmarshalMap(result.Item, &newItem); err != nil {
			return err
		}
	} else {
		newItem["shortKey"] = shortKey
		newItem[secondaryKeyName] = secondaryKeyValue
		newItem["counter"] = float64(0)
	}

	newItem["counter"] = newItem["counter"].(float64) + 1

	marshalNewItem, err := attributevalue.MarshalMap(newItem)

	if err != nil {
		return err
	}

	_, err = o.DbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: fullTableName,
		Item:      marshalNewItem,
	})

	if err != nil {
		return err
	}

	return nil
}

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
