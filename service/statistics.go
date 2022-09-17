package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	for _, err = range []error{
		o.IncrementStatisticsCounter("time_per_date", item.ShortKey, "datetime", RoundDateTimeAndConvertToTimestamp(&item.CreatedDate, 60*60*24)),
		o.IncrementStatisticsCounter("time_per_minute", item.ShortKey, "datetime", RoundDateTimeAndConvertToTimestamp(&item.CreatedDate, 60)),
		o.IncrementStatisticsCounter("referer", item.ShortKey, "referer", item.Referer),
		o.IncrementStatisticsCounter("devicetype", item.ShortKey, "devicetype", item.DeviceType),
	} {
		if err != nil {
			return err
		}
	}

	return nil
}

type StatisticsLegendType struct {
	TableName   string
	FieldName   string
	FieldTypeEx interface{}
}

type QueryStatisticsResultRow struct {
	Legend  interface{}
	Counter int
}

var (
	StatisticLegendTypeReferer       = &StatisticsLegendType{"referer", "referer", string("")}
	StatisticLegendTypeTimePerMinute = &StatisticsLegendType{"time_per_minute", "datetime", time.Time{}}
	StatisticLegendTypeTimePerDate   = &StatisticsLegendType{"time_per_date", "datetime", time.Time{}}
	StatisticLegendTypeDevicetype    = &StatisticsLegendType{"devicetype", "devicetype", string("")}
)

func (o *Service) QueryStatistics(shortKey string, legend *StatisticsLegendType, isAscending bool) ([]*QueryStatisticsResultRow, error) {
	fullTableName := aws.String("statistics_accum_" + legend.TableName)

	result, err := o.DbClient.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              fullTableName,
		KeyConditionExpression: aws.String("shortKey=:shortKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":shortKey": &types.AttributeValueMemberS{
				Value: shortKey,
			},
		},
		ScanIndexForward: &isAscending,
	})

	if err != nil {
		return nil, err
	}

	ret := make([]*QueryStatisticsResultRow, 0)
	decoder := attributevalue.NewDecoder()

	for _, item := range result.Items {
		legendReflect := reflect.New(reflect.ValueOf(legend.FieldTypeEx).Type())
		legendVal := legendReflect.Interface()
		err = decoder.Decode(item[legend.FieldName], legendVal)

		if err != nil {
			return nil, err
		}

		var counterVal int
		err = decoder.Decode(item["counter"], &counterVal)

		if err != nil {
			return nil, err
		}

		ret = append(ret, &QueryStatisticsResultRow{
			Counter: counterVal,
			Legend:  legendReflect.Elem().Interface(),
		})
	}

	return ret, nil
}

func RoundDateTimeAndConvertToTimestamp(target *time.Time, seconds int) int64 {
	secs := (int64(target.UnixMilli()) / 1000 / int64(seconds)) * int64(seconds)
	return secs
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
