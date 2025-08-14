package dynamodb

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestCreateLocalClient(t *testing.T) {
	tests := []struct {
		name string
		want *dynamodb.Client
	}{
		{
			name: "create local client",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, CreateLocalClient())
		})
	}
}
