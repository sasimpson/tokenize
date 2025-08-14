package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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
