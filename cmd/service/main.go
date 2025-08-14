package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"tokenize/api"
	"tokenize/persistence/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsdynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	//setup for dynamodb
	db := dynamodb.CreateLocalClient()

	_, err := db.DescribeTable(context.Background(), &awsdynamo.DescribeTableInput{
		TableName: aws.String("token_dat"),
	})
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			_ = dynamodb.CreateTable(context.Background(), db)
		}
	}

	handlers := &api.BaseHandler{
		Store: &dynamodb.DynamoStore{
			Api: db,
		},
	}

	routes := api.Routes(handlers)
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	fmt.Println("starting service on", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
