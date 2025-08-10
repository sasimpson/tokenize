package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToken_Encrypt(t *testing.T) {
	tests := []struct {
		name    string
		token   Token
		wantErr bool
	}{
		{
			name: "encrypt simple payload",
			token: Token{
				BaseModel: BaseModel{
					Id:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateToken: CreateToken{
					Payload:   "test payload",
					TokenType: "access",
					TTL:       3600,
					Metadata:  map[string]any{"key": "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "encrypt empty payload",
			token: Token{
				BaseModel: BaseModel{
					Id:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateToken: CreateToken{
					Payload:   "",
					TokenType: "refresh",
					TTL:       7200,
					Metadata:  map[string]any{},
				},
			},
			wantErr: false,
		},
		{
			name: "encrypt long payload",
			token: Token{
				BaseModel: BaseModel{
					Id:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateToken: CreateToken{
					Payload:   "this is a very long payload that contains lots of text and should still encrypt properly without any issues",
					TokenType: "session",
					TTL:       1800,
					Metadata:  map[string]any{"role": "admin", "permissions": []string{"read", "write"}},
				},
			},
			wantErr: false,
		},
		{
			name: "encrypt payload with special characters",
			token: Token{
				BaseModel: BaseModel{
					Id:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateToken: CreateToken{
					Payload:   "payload with special chars: !@#$%^&*(){}[]|\\:;\"'<>,.?/~`",
					TokenType: "api",
					TTL:       86400,
					Metadata:  map[string]any{"special": true},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPayload := tt.token.Payload

			err := tt.token.Encrypt()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEqual(t, originalPayload, tt.token.Payload, "payload should be encrypted")
			assert.NotEmpty(t, tt.token.Payload, "encrypted payload should not be empty")

			// Verify the encrypted payload is hex-encoded
			assert.Regexp(t, "^[0-9a-f]+$", tt.token.Payload, "encrypted payload should be hex-encoded")

			// Verify we can decrypt back to original
			decrypted, err := tt.token.Decrypt()
			assert.NoError(t, err)
			assert.Equal(t, originalPayload, decrypted, "decrypted payload should match original")
		})
	}
}

func TestToken_Decrypt(t *testing.T) {
	tests := []struct {
		name           string
		setupToken     func() Token
		expectedResult string
		wantErr        bool
	}{
		{
			name: "decrypt simple payload",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "test payload",
						TokenType: "access",
						TTL:       3600,
						Metadata:  map[string]any{"key": "value"},
					},
				}
				_ = token.Encrypt()
				return token
			},
			expectedResult: "test payload",
			wantErr:        false,
		},
		{
			name: "decrypt empty payload",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "",
						TokenType: "refresh",
						TTL:       7200,
						Metadata:  map[string]any{},
					},
				}
				_ = token.Encrypt()
				return token
			},
			expectedResult: "",
			wantErr:        false,
		},
		{
			name: "decrypt long payload",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "this is a very long payload that contains lots of text and should still decrypt properly without any issues",
						TokenType: "session",
						TTL:       1800,
						Metadata:  map[string]any{"role": "admin", "permissions": []string{"read", "write"}},
					},
				}
				_ = token.Encrypt()
				return token
			},
			expectedResult: "this is a very long payload that contains lots of text and should still decrypt properly without any issues",
			wantErr:        false,
		},
		{
			name: "decrypt payload with special characters",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "payload with special chars: !@#$%^&*(){}[]|\\:;\"'<>,.?/~`",
						TokenType: "api",
						TTL:       86400,
						Metadata:  map[string]any{"special": true},
					},
				}
				_ = token.Encrypt()
				return token
			},
			expectedResult: "payload with special chars: !@#$%^&*(){}[]|\\:;\"'<>,.?/~`",
			wantErr:        false,
		},
		{
			name: "decrypt invalid hex payload",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "invalid hex string",
						TokenType: "access",
						TTL:       3600,
						Metadata:  map[string]any{},
					},
				}
				return token
			},
			expectedResult: "invalid hex string",
			wantErr:        true,
		},
		{
			name: "decrypt corrupted ciphertext",
			setupToken: func() Token {
				token := Token{
					BaseModel: BaseModel{
						Id:        uuid.New(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CreateToken: CreateToken{
						Payload:   "deadbeef",
						TokenType: "access",
						TTL:       3600,
						Metadata:  map[string]any{},
					},
				}
				return token
			},
			expectedResult: "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()

			result, err := token.Decrypt()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestToken_Tokenize(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "tokenize simple payload",
			payload: "test payload",
			wantErr: false,
		},
		{
			name:    "tokenize empty payload",
			payload: "",
			wantErr: false,
		},
		{
			name:    "tokenize long payload",
			payload: "this is a very long payload that contains lots of text and should still generate a consistent hash",
			wantErr: false,
		},
		{
			name:    "tokenize payload with special characters",
			payload: "payload with special chars: !@#$%^&*(){}[]|\\:;\"'<>,.?/~`",
			wantErr: false,
		},
		{
			name:    "tokenize unicode payload",
			payload: "payload with unicode: 你好世界 🚀 ñáéíóú",
			wantErr: false,
		},
		{
			name:    "tokenize json payload",
			payload: `{"user":"admin","role":"super","permissions":["read","write","delete"]}`,
			wantErr: false,
		},
		{
			name:    "tokenize whitespace payload",
			payload: "   \t\n   ",
			wantErr: false,
		},
		{
			name:    "tokenize numeric string payload",
			payload: "1234567890",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := Token{
				BaseModel: BaseModel{
					Id:        uuid.New(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				CreateToken: CreateToken{
					Payload:   tt.payload,
					TokenType: "test",
					TTL:       3600,
					Metadata:  map[string]any{},
				},
			}

			err := token.Tokenize()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, token.Token, "token should not be empty")
			assert.Len(t, token.Token, 64, "SHA512-256 hash should be 64 hex characters")
			assert.Regexp(t, "^[0-9a-f]+$", token.Token, "token should be hex-encoded")

			// Test deterministic behavior - same payload should generate same token
			token2 := Token{
				CreateToken: CreateToken{
					Payload: tt.payload,
				},
			}
			err2 := token2.Tokenize()
			assert.NoError(t, err2)
			assert.Equal(t, token.Token, token2.Token, "same payload should generate same token")
		})
	}
}
