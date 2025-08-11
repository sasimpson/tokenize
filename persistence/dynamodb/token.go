package dynamodb

import (
	"context"
	"time"

	"tokenize/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

func (d *DynamoStore) GetToken(ctx context.Context, token string) (*models.Token, error) {
	awsTokenVal, err := attributevalue.Marshal(token)
	if err != nil {
		return nil, err
	}

	dynamoItem, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: TokenTableName,
		Key: map[string]types.AttributeValue{
			"token": awsTokenVal,
		},
	})
	if err != nil {
		return nil, err
	}

	tokenPayload := &models.Token{}
	err = attributevalue.UnmarshalMap(dynamoItem.Item, tokenPayload)
	if err != nil {
		return nil, err
	}
	return tokenPayload, nil
}

func (d *DynamoStore) CreateToken(ctx context.Context, token *models.Token) (*models.Token, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	token.Id = id
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	dynamoItem, err := attributevalue.MarshalMap(token)
	if err != nil {
		return nil, err
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: TokenTableName,
		Item:      dynamoItem,
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// func (d *DynamoStore) UpdateToken(ctx context.Context, token *models.Token) (*models.Token, error) {}
func (d *DynamoStore) DeleteToken(ctx context.Context, token *models.Token) error {
	awsTokenVal, err := attributevalue.Marshal(token)
	if err != nil {
		return err
	}
	_, err = d.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: TokenTableName,
		Key: map[string]types.AttributeValue{
			"token": awsTokenVal,
		},
	})

	return err
}
