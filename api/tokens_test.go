package api

import (
	"context"
	"errors"
	"testing"

	"tokenize/models"
	"tokenize/persistence"
	"tokenize/persistence/mock"

	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateToken(t *testing.T) {
	type fields struct {
		Store persistence.Store
	}
	type args struct {
		ctx context.Context
		in  *NewTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Create a token",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token:     "e3061477f33275654a7beebe7ac6a4941adedec434d870a8bac71e7bff2eb137",
						BaseModel: models.BaseModel{},
						CreateToken: models.CreateToken{
							Payload:   "fc8df3ea16c7823811c85fead07f6589684f9799084fbad080134cc16d512339f5dfe9",
							TTL:       3600,
							TokenType: "access",
							Metadata:  map[string]any{"foo": "bar"},
						},
					},
				},
			},
			args: args{
				ctx: context.Background(),
				in: &NewTokenRequest{
					Body: struct {
						Data models.CreateToken `json:"data" validate:"required"`
					}{
						Data: models.CreateToken{
							Payload:   "this is the payload",
							TTL:       3600,
							TokenType: "access",
							Metadata:  map[string]any{"foo": "bar"},
						},
					},
				},
			},
			want: "e3061477f33275654a7beebe7ac6a4941adedec434d870a8bac71e7bff2eb137",
		}, {
			name: "create token error",
			fields: fields{
				Store: mock.Store{
					CreateError: errors.New("unknown error"),
				},
			},
			args: args{
				ctx: context.Background(),
				in: &NewTokenRequest{
					Body: struct {
						Data models.CreateToken `json:"data" validate:"required"`
					}{
						Data: models.CreateToken{
							Payload:   "this is the payload",
							TTL:       3600,
							TokenType: "access",
							Metadata:  map[string]any{"foo": "bar"},
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				Store: tt.fields.Store,
			}
			got, err := h.CreateToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got.Body.Token)
		})
	}
}

func TestHandler_GetEncryptedToken(t *testing.T) {
	type fields struct {
		Store persistence.Store
	}
	type args struct {
		ctx context.Context
		in  *GetTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Token
		wantErr bool
	}{
		{
			name: "Get a token",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token:     "foobartesttoken",
						BaseModel: models.BaseModel{},
						CreateToken: models.CreateToken{
							Payload: "this is the payload",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: "foobartesttoken"},
			},
			want: &models.Token{
				Token:     "foobartesttoken",
				BaseModel: models.BaseModel{},
				CreateToken: models.CreateToken{
					Payload: "",
				},
			},
		}, {
			name: "no token for getting",
			fields: fields{
				Store: mock.Store{
					GetError: models.ErrTokenNotFound,
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: ""},
			},
			wantErr: true,
		}, {
			name: "no token returned",
			fields: fields{
				Store: mock.Store{
					GetError: models.ErrTokenNotFound,
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: "foobartesttoken"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				Store: tt.fields.Store,
			}
			got, err := h.GetEncryptedToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, *tt.want, got.Body.Token)
		})
	}
}

func TestHandler_GetDecryptedToken(t *testing.T) {
	type fields struct {
		Store persistence.Store
	}
	type args struct {
		ctx context.Context
		in  *GetTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Token
		wantErr bool
	}{
		{
			name: "Get a token",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token:     "foobartesttoken",
						BaseModel: models.BaseModel{},
						CreateToken: models.CreateToken{
							Payload: "fc8df3ea16c7823811c85fead07f6589684f9799084fbad080134cc16d512339f5dfe9",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: "foobartesttoken"},
			},
			want: &models.Token{
				Token:     "foobartesttoken",
				BaseModel: models.BaseModel{},
				CreateToken: models.CreateToken{
					Payload: "this is the payload",
				},
			},
		}, {
			name: "No token for getting",
			fields: fields{
				Store: mock.Store{},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: ""},
			},
			wantErr: true,
		}, {
			name: "no token returned",
			fields: fields{
				Store: mock.Store{
					GetError: models.ErrTokenNotFound,
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: "foobartesttoken"},
			},
			wantErr: true,
		}, {
			name: "Decrypt error",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token:     "foobartesttoken",
						BaseModel: models.BaseModel{},
						CreateToken: models.CreateToken{
							Payload: "deadbeef",
						},
					},
				},
			},
			args: args{
				ctx: context.Background(), in: &GetTokenRequest{Token: "foobartesttoken"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				Store: tt.fields.Store,
			}
			got, err := h.GetDecryptedToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, *tt.want, got.Body.Token)
		})
	}
}

func TestHandler_DeleteToken(t *testing.T) {
	type fields struct {
		Store persistence.Store
	}
	type args struct {
		ctx context.Context
		in  *GetTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *struct{}
		wantErr bool
	}{
		{
			name: "Delete a token",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token: "foobartesttoken",
					},
				},
			},
			args: args{
				ctx: context.Background(),
				in: &GetTokenRequest{
					Token: "foobartesttoken",
				},
			},
		}, {
			name: "No token to delete",
			fields: fields{
				Store: mock.Store{
					GetError: models.ErrTokenNotFound,
				},
			},
			args: args{
				ctx: context.Background(),
				in: &GetTokenRequest{
					Token: "foobartesttoken",
				},
			},
			wantErr: true,
		}, {
			name: "token delete error",
			fields: fields{
				Store: mock.Store{
					Token: &models.Token{
						Token: "foobartesttoken",
					},
					DeleteError: errors.New("unknown error"),
				},
			},
			args: args{
				ctx: context.Background(),
				in: &GetTokenRequest{
					Token: "foobartesttoken",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BaseHandler{
				Store: tt.fields.Store,
			}
			got, err := h.DeleteToken(tt.args.ctx, tt.args.in)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
