package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"tokenize/api"
	"tokenize/persistence/dynamodb"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsdynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	//setup for dynamodb
	awsregion := os.Getenv("AWS_REGION")
	if awsregion == "" {
		awsregion = "localhost"
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsregion))
	if err != nil {
		panic(err)
	}

	db := awsdynamo.NewFromConfig(cfg, func(o *awsdynamo.Options) {
		o.BaseEndpoint = aws.String("http://localhost:8000")
	})

	_, err = db.DescribeTable(context.Background(), &awsdynamo.DescribeTableInput{
		TableName: aws.String("token_dat"),
	})
	if err != nil {
		var notFoundEx *types.ResourceNotFoundException
		if errors.As(err, &notFoundEx) {
			dynamodb.CreateTable(context.Background(), db)
		}
	}

	handlers := &api.BaseHandler{
		Store: &dynamodb.DynamoStore{
			Client: db,
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
