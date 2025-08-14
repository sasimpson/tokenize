package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func SetupDynamoTable(ctx context.Context, client Api) {
	_, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String("token_data"),
	})
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			_ = CreateTable(ctx, client)
		}
	}
	return
}

func CreateTable(ctx context.Context, client Api) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: TokenTableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("token"),
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
	return err
}
