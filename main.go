package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/borowiak-m/go-microservice/routes"
	"github.com/borowiak-m/go-microservice/server"
	"github.com/gorilla/mux"
)

func main() {
	// Create instances for router, logger, validation
	muxRouter := mux.NewRouter()
	newlogger := log.New(os.Stdout, "product-api", log.LstdFlags)
	newValidation := data.NewValidation()
	// SERVER
	server := server.RunServer(muxRouter)
	// HANDLERS
	routes.RegisterProductRoutes(newlogger, newValidation, muxRouter)
	routes.RegisterUserRoutes(newlogger, newValidation, muxRouter)
	//
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// waiting on main process for termination signals to greacefully handle the event
	// potentially close any db connections etc
	sig := <-sigChan // blocking until receives a signal from sigChan
	log.Println("Received terminate, greaceful shutdown", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
