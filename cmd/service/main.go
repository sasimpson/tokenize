package main

import (
	"context"
	"fmt"
	"net/http"

	"tokenize/api"
	"tokenize/persistence/dynamodb"
)

func main() {
	db := dynamodb.CreateLocalClient()

	dynamodb.SetupDynamoTable(context.Background(), db)
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
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
