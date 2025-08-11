package api

import (
	"context"
	"testing"

	"tokenize/models"
	"tokenize/persistence"
	"tokenize/persistence/mock"

	"github.com/stretchr/testify/assert"
)

func TestBaseHandler_DeleteToken(t *testing.T) {
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
					GetError:    nil,
					DeleteError: nil,
				},
			},
			args: args{
				ctx: context.Background(),
				in: &GetTokenRequest{
					Token: "foobartesttoken",
				},
			},
			want:    nil,
			wantErr: false,
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

func TestBaseHandler_GetEncryptedToken(t *testing.T) {
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

func TestBaseHandler_GetDecryptedToken(t *testing.T) {
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
