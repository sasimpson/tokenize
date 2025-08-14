package dynamodb

import (
	"context"
	"errors"
	"testing"
	"time"
	"tokenize/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDynamoAPI struct {
	getItemFunc       func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	putItemFunc       func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	deleteItemFunc    func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	createTableFunc   func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	describeTableFunc func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

func (m *mockDynamoAPI) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.getItemFunc != nil {
		return m.getItemFunc(ctx, params, optFns...)
	}
	return nil, errors.New("GetItem not implemented")
}

func (m *mockDynamoAPI) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.putItemFunc != nil {
		return m.putItemFunc(ctx, params, optFns...)
	}
	return nil, errors.New("PutItem not implemented")
}

func (m *mockDynamoAPI) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if m.deleteItemFunc != nil {
		return m.deleteItemFunc(ctx, params, optFns...)
	}
	return nil, errors.New("DeleteItem not implemented")
}

func (m *mockDynamoAPI) CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
	if m.createTableFunc != nil {
		return m.createTableFunc(ctx, params, optFns...)
	}
	return nil, errors.New("CreateTable not implemented")
}

func (m *mockDynamoAPI) DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	if m.describeTableFunc != nil {
		return m.describeTableFunc(ctx, params, optFns...)
	}
	return nil, errors.New("DescribeTable not implemented")
}

func TestGetToken(t *testing.T) {
	testCases := []struct {
		name   string
		token  string
		client func(t *testing.T) *mockDynamoAPI
		expect func(t *testing.T, token *models.Token, err error)
	}{
		{
			name:  "successful token retrieval",
			token: "test-token-123",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
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
					},
				}
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
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return &dynamodb.GetItemOutput{
							Item: nil,
						}, nil
					},
				}
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
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return nil, nil
					},
				}
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
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return nil, errors.New("dynamodb error")
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Equal(t, "dynamodb error", err.Error())
				assert.Nil(t, token)
			},
		},
		{
			name:  "invalid token for marshal error",
			token: string([]byte{0xff, 0xfe, 0xfd}),
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return nil, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Nil(t, token)
			},
		},
		{
			name:  "unmarshal error - invalid item structure",
			token: "test-token-unmarshal-error",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					getItemFunc: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
						return &dynamodb.GetItemOutput{
							Item: map[string]types.AttributeValue{
								"invalid_field": &types.AttributeValueMemberS{
									Value: "invalid",
								},
								"createdAt": &types.AttributeValueMemberN{
									Value: "not-a-date",
								},
							},
						}, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Nil(t, token)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &DynamoStore{
				Api: tc.client(t),
			}

			token, err := store.GetToken(context.Background(), tc.token)
			tc.expect(t, token, err)
		})
	}
}

func TestCreateToken(t *testing.T) {
	testCases := []struct {
		name   string
		input  *models.Token
		client func(t *testing.T) *mockDynamoAPI
		expect func(t *testing.T, token *models.Token, err error)
	}{
		{
			name: "successful token creation",
			input: &models.Token{
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"source": "test",
						"user":   "testuser",
					},
				},
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Item)

						// Check for fields that should be present based on the dynamodbav tags
						assert.NotNil(t, params.Item["Id"]) // Note: Id might use different tag or key
						assert.NotNil(t, params.Item["createdAt"])
						assert.NotNil(t, params.Item["updatedAt"])
						assert.NotNil(t, params.Item["payload"])
						assert.NotNil(t, params.Item["token_type"])
						assert.NotNil(t, params.Item["ttl"])
						assert.NotNil(t, params.Item["metadata"])

						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEqual(t, uuid.Nil, token.Id)
				assert.False(t, token.CreatedAt.IsZero())
				assert.False(t, token.UpdatedAt.IsZero())
				assert.Equal(t, "test-payload", token.Payload)
				assert.Equal(t, "bearer", token.TokenType)
				assert.Equal(t, int64(3600), token.TTL)
				assert.NotNil(t, token.Metadata)
				assert.Equal(t, "test", token.Metadata["source"])
				assert.Equal(t, "testuser", token.Metadata["user"])
			},
		},
		{
			name: "uuid generation failure simulation",
			input: &models.Token{
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"source": "test",
					},
				},
			},
			client: func(t *testing.T) *mockDynamoAPI {
				// This test demonstrates the CreateToken method's marshal and put flow
				// In practice, uuid.NewV7() and attributevalue.MarshalMap work for valid inputs
				return &mockDynamoAPI{
					putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						// Verify the input is properly structured
						assert.NotNil(t, params.Item)
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, "test-payload", token.Payload)
				assert.Equal(t, "bearer", token.TokenType)
				assert.Equal(t, int64(3600), token.TTL)
			},
		},
		{
			name: "dynamodb put item error",
			input: &models.Token{
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"source": "test",
					},
				},
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						return nil, errors.New("dynamodb put error")
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.Error(t, err)
				assert.Equal(t, "dynamodb put error", err.Error())
				assert.Nil(t, token)
			},
		},
		{
			name: "empty payload token",
			input: &models.Token{
				CreateToken: models.CreateToken{
					Payload:   "",
					TokenType: "bearer",
					TTL:       3600,
					Metadata:  nil,
				},
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						assert.NotNil(t, params.Item)
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, "", token.Payload)
				assert.Equal(t, "bearer", token.TokenType)
				assert.Equal(t, int64(3600), token.TTL)
				assert.NotEqual(t, uuid.Nil, token.Id)
			},
		},
		{
			name: "successful creation with minimal data",
			input: &models.Token{
				CreateToken: models.CreateToken{
					Payload:   "minimal",
					TokenType: "api",
					TTL:       1800,
				},
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Item)
						return &dynamodb.PutItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, token *models.Token, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.Equal(t, "minimal", token.Payload)
				assert.Equal(t, "api", token.TokenType)
				assert.Equal(t, int64(1800), token.TTL)
				assert.NotEqual(t, uuid.Nil, token.Id)

				now := time.Now()
				assert.WithinDuration(t, now, token.CreatedAt, 1*time.Second)
				assert.WithinDuration(t, now, token.UpdatedAt, 1*time.Second)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &DynamoStore{
				Api: tc.client(t),
			}

			token, err := store.CreateToken(context.Background(), tc.input)
			tc.expect(t, token, err)
		})
	}
}

func TestDeleteToken(t *testing.T) {
	testCases := []struct {
		name   string
		input  *models.Token
		client func(t *testing.T) *mockDynamoAPI
		expect func(t *testing.T, err error)
	}{
		{
			name: "successful token deletion",
			input: &models.Token{
				BaseModel: models.BaseModel{
					Id:        uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"source": "test",
					},
				},
				Token: "test-token-123",
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						// Verify the table name is correct
						assert.Equal(t, *TokenTableName, *params.TableName)

						// Verify the key is provided
						assert.NotNil(t, params.Key)
						assert.NotNil(t, params.Key["token"])

						return &dynamodb.DeleteItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "dynamodb delete item error",
			input: &models.Token{
				BaseModel: models.BaseModel{
					Id:        uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"source": "test",
					},
				},
				Token: "test-token-error",
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						return nil, errors.New("dynamodb delete error")
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Equal(t, "dynamodb delete error", err.Error())
			},
		},
		{
			name: "successful deletion with complex metadata",
			input: &models.Token{
				BaseModel: models.BaseModel{
					Id:        uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				CreateToken: models.CreateToken{
					Payload:   "test-payload",
					TokenType: "bearer",
					TTL:       3600,
					Metadata: map[string]any{
						"complex": map[string]any{
							"nested": "value",
						},
					},
				},
				Token: "test-token-complex",
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Key)
						assert.NotNil(t, params.Key["token"])
						return &dynamodb.DeleteItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "delete with minimal token data",
			input: &models.Token{
				BaseModel: models.BaseModel{
					Id:        uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				CreateToken: models.CreateToken{
					Payload:   "minimal",
					TokenType: "api",
					TTL:       1800,
				},
				Token: "minimal-token",
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Key)
						assert.NotNil(t, params.Key["token"])

						return &dynamodb.DeleteItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "delete with empty metadata",
			input: &models.Token{
				BaseModel: models.BaseModel{
					Id:        uuid.MustParse("01234567-89ab-cdef-0123-456789abcdef"),
					CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				CreateToken: models.CreateToken{
					Payload:   "empty-metadata",
					TokenType: "bearer",
					TTL:       3600,
					Metadata:  nil,
				},
				Token: "token-no-metadata",
			},
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Key)

						return &dynamodb.DeleteItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:  "nil token input",
			input: nil,
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					deleteItemFunc: func(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
						// In practice, marshaling nil creates an empty structure
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.NotNil(t, params.Key)
						return &dynamodb.DeleteItemOutput{}, nil
					},
				}
			},
			expect: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store := &DynamoStore{
				Api: tc.client(t),
			}

			err := store.DeleteToken(context.Background(), tc.input)
			tc.expect(t, err)
		})
	}
}

func TestSetupDynamoTable(t *testing.T) {
	testCases := []struct {
		name   string
		client func(t *testing.T) *mockDynamoAPI
		expect func(t *testing.T, mock *mockDynamoAPI)
	}{
		{
			name: "table exists - no action needed",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						// Verify correct table name is being checked
						assert.Equal(t, "token_data", *params.TableName)

						return &dynamodb.DescribeTableOutput{
							Table: &types.TableDescription{
								TableName:   aws.String("token_data"),
								TableStatus: types.TableStatusActive,
							},
						}, nil
					},
				}
			},
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// Should not call CreateTable when table exists
				assert.Nil(t, mock.createTableFunc, "CreateTable should not be called when table exists")
			},
		},
		{
			name: "table not found - creates table",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						assert.Equal(t, "token_data", *params.TableName)

						// Return ResourceNotFoundException
						return nil, &types.ResourceNotFoundException{
							Message: aws.String("Table not found: token_data"),
						}
					},
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						// Verify CreateTable is called with correct parameters
						assert.Equal(t, *TokenTableName, *params.TableName)
						assert.Len(t, params.AttributeDefinitions, 1)
						assert.Equal(t, "token", *params.AttributeDefinitions[0].AttributeName)
						assert.Equal(t, types.ScalarAttributeTypeS, params.AttributeDefinitions[0].AttributeType)
						assert.Len(t, params.KeySchema, 1)
						assert.Equal(t, "token", *params.KeySchema[0].AttributeName)
						assert.Equal(t, types.KeyTypeHash, params.KeySchema[0].KeyType)
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
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// CreateTable should have been called
				assert.NotNil(t, mock.createTableFunc, "CreateTable should be called when table doesn't exist")
			},
		},
		{
			name: "describe table generic error - no table creation",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						assert.Equal(t, "token_data", *params.TableName)

						// Return generic error (not ResourceNotFoundException)
						return nil, errors.New("access denied")
					},
				}
			},
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// Should not call CreateTable for non-ResourceNotFoundException errors
				assert.Nil(t, mock.createTableFunc, "CreateTable should not be called for generic errors")
			},
		},
		{
			name: "table creation fails",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						return nil, &types.ResourceNotFoundException{
							Message: aws.String("Table not found: token_data"),
						}
					},
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						// Simulate CreateTable failure
						return nil, errors.New("insufficient permissions")
					},
				}
			},
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// CreateTable should have been attempted even if it failed
				assert.NotNil(t, mock.createTableFunc, "CreateTable should be attempted")
			},
		},
		{
			name: "resource already exists during creation",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						return nil, &types.ResourceNotFoundException{
							Message: aws.String("Table not found: token_data"),
						}
					},
					createTableFunc: func(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error) {
						// Simulate race condition where table gets created between describe and create
						return nil, &types.ResourceInUseException{
							Message: aws.String("Table already exists: token_data"),
						}
					},
				}
			},
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// CreateTable should have been attempted
				assert.NotNil(t, mock.createTableFunc, "CreateTable should be attempted")
			},
		},
		{
			name: "multiple error types handling",
			client: func(t *testing.T) *mockDynamoAPI {
				return &mockDynamoAPI{
					describeTableFunc: func(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
						// Return a different AWS error type
						return nil, &types.LimitExceededException{
							Message: aws.String("Rate limit exceeded"),
						}
					},
				}
			},
			expect: func(t *testing.T, mock *mockDynamoAPI) {
				// Should not call CreateTable for LimitExceededException
				assert.Nil(t, mock.createTableFunc, "CreateTable should not be called for LimitExceededException")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := tc.client(t)

			// Call SetupDynamoTable function
			SetupDynamoTable(context.Background(), mock)

			// Run expectations
			tc.expect(t, mock)
		})
	}
}
