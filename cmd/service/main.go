package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tokenize/api"
	"tokenize/persistence/dynamodb"
)

func buildServer() *http.Server {
	db := dynamodb.CreateLocalClient()

	dynamodb.SetupDynamoTable(context.Background(), db)
	handlers := &api.BaseHandler{
		Store: &dynamodb.DynamoStore{
			Api: db,
		},
	}
	routes := api.Routes(handlers)
	return &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

}

func main() {
	server := buildServer()

	shutdownChan := make(chan bool, 1)

	fmt.Println("starting service on", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v\n", err)
		}
		time.Sleep(1 * time.Millisecond)
		log.Println("stopped serving new connections")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10+time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("http shutdown error: %v\n", err)
	}

	<-shutdownChan
	log.Println("graceful shutdown complete")
}
