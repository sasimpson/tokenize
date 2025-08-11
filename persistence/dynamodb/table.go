package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateTable(ctx context.Context, client *dynamodb.Client) {
	client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: TokenTableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("token"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("type"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("ttl"),
				AttributeType: types.ScalarAttributeTypeN,
			}, {
				AttributeName: aws.String("created_at"),
				AttributeType: types.ScalarAttributeTypeN,
			}, {
				AttributeName: aws.String("updated_at"),
				AttributeType: types.ScalarAttributeTypeN,
			}, {
				AttributeName: aws.String("metadata"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("payload"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("token"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
}
