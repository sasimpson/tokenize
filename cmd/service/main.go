package main

import (
	"context"
	"errors"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	server := buildServer()

	shutdownChan := make(chan bool, 1)

	slog.Info("starting service", "addr", server.Addr)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
		time.Sleep(1 * time.Millisecond)
		slog.Info("stopped serving new connections")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10+time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("http shutdown error", "error", err)
		os.Exit(1)
	}

	<-shutdownChan
	slog.Info("graceful shutdown complete")
}
