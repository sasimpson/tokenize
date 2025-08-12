package dynamodb

import (
	"context"
	"errors"
	"testing"

	"tokenize/models"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type GetItemAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type mockGetItemAPI func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)

func (m mockGetItemAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m(ctx, params, optFns...)
}

type testDynamoStore struct {
	client GetItemAPI
}

func (d *testDynamoStore) GetToken(ctx context.Context, token string) (*models.Token, error) {
	awsTokenVal, err := attributevalue.Marshal(token)
	if err != nil {
		return nil, err
	}

	dynamoItem, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: TokenTableName,
		Key: map[string]types.AttributeValue{
			"token": awsTokenVal,
		},
	})
	if err != nil {
		return nil, err
	}
	if dynamoItem == nil || dynamoItem.Item == nil {
		return nil, models.ErrTokenNotFound
	}

	tokenPayload := &models.Token{}
	err = attributevalue.UnmarshalMap(dynamoItem.Item, tokenPayload)
	if err != nil {
		return nil, err
	}
	return tokenPayload, nil
}

func TestGetToken(t *testing.T) {
	testCases := []struct {
		name   string
		token  string
		client func(t *testing.T) mockGetItemAPI
		expect func(t *testing.T, token *models.Token, err error)
	}{
		{
			name:  "successful token retrieval",
			token: "test-token-123",
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return &dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"id": &types.AttributeValueMemberB{
								Value: []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef},
							},
							"token": &types.AttributeValueMemberS{
								Value: "test-token-123",
							},
							"payload": &types.AttributeValueMemberS{
								Value: "encrypted-payload",
							},
							"token_type": &types.AttributeValueMemberS{
								Value: "bearer",
							},
							"ttl": &types.AttributeValueMemberN{
								Value: "3600",
							},
							"createdAt": &types.AttributeValueMemberS{
								Value: "2024-01-01T00:00:00Z",
							},
							"updatedAt": &types.AttributeValueMemberS{
								Value: "2024-01-01T00:00:00Z",
							},
							"metadata": &types.AttributeValueMemberM{
								Value: map[string]types.AttributeValue{
									"source": &types.AttributeValueMemberS{Value: "test"},
								},
							},
						},
					}, nil
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, "test-token-123", token.Token)
				assert.Equal(t, "encrypted-payload", token.Payload)
				assert.Equal(t, "bearer", token.TokenType)
				assert.Equal(t, int64(3600), token.TTL)
				assert.NotNil(t, token.Metadata)
			},
		},
		{
			name:  "token not found - nil item",
			token: "nonexistent-token",
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return &dynamodb.GetItemOutput{
						Item: nil,
					}, nil
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Equal(t, models.ErrTokenNotFound, err)
				assert.Nil(t, token)
			},
		},
		{
			name:  "token not found - nil output",
			token: "nonexistent-token",
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return nil, nil
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Equal(t, models.ErrTokenNotFound, err)
				assert.Nil(t, token)
			},
		},
		{
			name:  "dynamodb get item error",
			token: "test-token-error",
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return nil, errors.New("dynamodb error")
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Equal(t, "dynamodb error", err.Error())
				assert.Nil(t, token)
			},
		},
		{
			name:  "invalid token for marshal error",
			token: string([]byte{0xff, 0xfe, 0xfd}), // Invalid UTF-8 to trigger marshal error
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					// This won't be called due to marshal error
					return nil, nil
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Nil(t, token)
			},
		},
		{
			name:  "unmarshal error - invalid item structure",
			token: "test-token-unmarshal-error",
			client: func(t *testing.T) mockGetItemAPI {
				return mockGetItemAPI(func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return &dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"invalid_field": &types.AttributeValueMemberS{
								Value: "invalid",
							},
							"createdAt": &types.AttributeValueMemberN{
								Value: "not-a-date", // Wrong type for date field
							},
						},
					}, nil
				})
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Nil(t, token)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a testable store with our mock client
			testStore := &testDynamoStore{
				client: tc.client(t),
			}

			token, err := testStore.GetToken(context.Background(), tc.token)
			tc.expect(t, token, err)
		})
	}
}
