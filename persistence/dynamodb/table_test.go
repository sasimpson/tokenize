package dynamodb

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateTable(t *testing.T) {
	testCases := []struct {
		name   string
		client func(t *testing.T) *mockDynamoAPI
		expect func(t *testing.T, err error)
	}{
		{
			name: "successful table creation",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						// Verify table name
						assert.Equal(t, *TokenTableName, *params.TableName)

						// Verify attribute definitions
						assert.Len(t, params.AttributeDefinitions, 1)
						assert.Equal(t, "token", *params.AttributeDefinitions[0].AttributeName)
						assert.Equal(t, types.ScalarAttributeTypeS, params.AttributeDefinitions[0].AttributeType)

						// Verify key schema
						assert.Len(t, params.KeySchema, 1)
						assert.Equal(t, "token", *params.KeySchema[0].AttributeName)
						assert.Equal(t, types.KeyTypeHash, params.KeySchema[0].KeyType)

						// Verify billing mode
						assert.Equal(t, types.BillingModePayPerRequest, params.BillingMode)

						return &dynamodb.CreateTableOutput{
							TableDescription: &types.TableDescription{
								TableName:   params.TableName,
								TableStatus: types.TableStatusCreating,
							},
						}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "dynamodb create table error",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						return nil, errors.New("table already exists")
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "table already exists", err.Error())
			},
		},
		{
			name: "resource already exists error",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						return nil, &types.ResourceInUseException{
							Message: aws.String("Table already exists: token_data"),
						}
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
				var resourceInUseErr *types.ResourceInUseException
				assert.True(t, errors.As(err, &resourceInUseErr))
			},
		},
		{
			name: "verify table configuration parameters",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						// Detailed verification of all table parameters
						assert.Equal(t, "token_data", *params.TableName)

						// Check attribute definitions structure
						assert.Len(t, params.AttributeDefinitions, 1)
						attr := params.AttributeDefinitions[0]
						assert.Equal(t, "token", *attr.AttributeName)
						assert.Equal(t, types.ScalarAttributeTypeS, attr.AttributeType)

						// Check key schema structure
						assert.Len(t, params.KeySchema, 1)
						key := params.KeySchema[0]
						assert.Equal(t, "token", *key.AttributeName)
						assert.Equal(t, types.KeyTypeHash, key.KeyType)

						// Verify billing mode is pay-per-request
						assert.Equal(t, types.BillingModePayPerRequest, params.BillingMode)

						// Verify no provisioned throughput is set (since we're using pay-per-request)
						assert.Nil(t, params.ProvisionedThroughput)

						return &dynamodb.CreateTableOutput{
							TableDescription: &types.TableDescription{
								TableName:            params.TableName,
								TableStatus:          types.TableStatusActive,
								AttributeDefinitions: params.AttributeDefinitions,
								KeySchema:            params.KeySchema,
								BillingModeSummary: &types.BillingModeSummary{
									BillingMode: types.BillingModePayPerRequest,
								},
							},
						}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "limit exceeded error",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						return nil, &types.LimitExceededException{
							Message: aws.String("Too many tables in account"),
						}
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
				var limitExceededErr *types.LimitExceededException
				assert.True(t, errors.As(err, &limitExceededErr))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Now we can test the actual CreateTable function using our mock
			mock := tc.client(t)

			// Call the actual CreateTable function from table.go
			err := CreateTable(context.Background(), mock)

			tc.expect(t, err)
		})
	}
}
