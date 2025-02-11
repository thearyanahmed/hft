package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/config"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/handler"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/logger"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/repository"
	"github.com/hellofreshdevtests/HFtest-platform-engineering-thearyanahmed/pkg/service"
)

func main() {
	conf, err := config.FromENV()

	if err != nil {
		log.Fatal(err)
	}

	logger.SetupLogger(conf)
	svcLogger := logger.Logger()

	repo := repository.NewInMemoryRepository()

	svc := service.NewService(repo)

	router := handler.NewRouter(svc)

	address := conf.AppAddress()

	server := &http.Server{Addr: address, Handler: router}

	serverCtx, cancelFunc := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		// @todo should be configurable from .env
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		cancelFunc()
	}()

	// Run the server
	svcLogger.Printf("will be serving on: %s", address)

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
