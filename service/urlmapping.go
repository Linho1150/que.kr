package service

import (
	"context"
	"errors"
	"quekr/server/util"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type MappingInfo struct {
	ShortKey         string `dynamodbav:"innerId"`
	OriginalUrl      string
	CreatedDate      time.Time
	RequestIPAddress string
	SecretToken      string
}

var URLMAP_TABLENAME = "url_map"

func (o *Service) tryPersistMapping(mappingInfo *MappingInfo, allowOverwrite bool) error {
	marshal, err := attributevalue.MarshalMap(mappingInfo)

	if err != nil {
		return err
	}

	var conditionExpression *string

	if allowOverwrite {
		conditionExpression = aws.String("attribute_not_exists(Id)")
	}

	result, err := o.DbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           &URLMAP_TABLENAME,
		Item:                marshal,
		ConditionExpression: conditionExpression,
		ReturnValues:        types.ReturnValueAllOld,
	})

	if err != nil {
		return err
	}

	if _, exist := result.Attributes["Id"]; exist && !allowOverwrite {
		return errors.New("oops, its duplicated")
	}

	return nil
}

func (o *Service) tryCreateMapping(originalUrl string, requestIPAddress string) (*MappingInfo, error) {
	shortKey, err := util.GenerateShortKey()

	if err != nil {
		return nil, err
	}

	secretToken, err := util.GenerateSecretToken()

	if err != nil {
		return nil, err
	}

	mappingInfo := &MappingInfo{
		ShortKey:         shortKey,
		OriginalUrl:      originalUrl,
		CreatedDate:      time.Now(),
		SecretToken:      secretToken,
		RequestIPAddress: requestIPAddress,
	}

	if err = o.tryPersistMapping(mappingInfo, false); err != nil {
		return nil, err
	}

	return mappingInfo, nil
}

func (o *Service) CreateMapping(originalUrl string, requestIPAddress string) (*MappingInfo, error) {
	var lastErr error

	for trying := 0; trying < 3; trying++ {
		ret, err := o.tryCreateMapping(originalUrl, requestIPAddress)

		if err == nil {
			return ret, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

func (o *Service) RemoveMapping(shortKey string, secretToken string) error {
	item, err := o.QueryMapping(shortKey)

	if err != nil {
		return err
	}

	if item.SecretToken != secretToken {
		return errors.New("supplied secretToken invalid")
	}

	_, err = o.DbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			"innerId": &types.AttributeValueMemberS{
				Value: shortKey,
			},
		},
		TableName: &URLMAP_TABLENAME,
	})

	if err != nil {
		return err
	}

	return nil
}

func (o *Service) UpdateMapping(shortKey string, secretToken string, originalUrl string) error {
	item, err := o.QueryMapping(shortKey)

	if err != nil {
		return err
	}

	if item.SecretToken != secretToken {
		return errors.New("supplied secretToken invalid")
	}

	item.OriginalUrl = originalUrl
	err = o.tryPersistMapping(item, true)

	if err != nil {
		return err
	}

	return nil
}

func (o *Service) QueryMapping(shortKey string) (*MappingInfo, error) {
	result, err := o.DbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"innerId": &types.AttributeValueMemberS{
				Value: shortKey,
			},
		},
		TableName: &URLMAP_TABLENAME,
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("item does not exist")
	}


	item := &MappingInfo{}

	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, err
	}

	return item, nil
}
