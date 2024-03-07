package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/borowiak-m/go-microservice/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// new logger
	newlogger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// HANDLERS
	//   create handler for products with logger
	handlderProducts := handlers.NewProducts(newlogger)

	// SERVER
	//   new serve mux
	servMx := mux.NewRouter()
	//   register handlerProduct as server for "/products" pattern
	getRouter := servMx.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", handlderProducts.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.GetSingleProduct)

	putRouter := servMx.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.UpdateSingleProduct)
	// executes middleware before it can go to the HandleFunc
	putRouter.Use(handlderProducts.MiddlewareProductValidation)

	postRouter := servMx.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", handlderProducts.CreateProduct)
	// executes middleware before it can go to the HandleFunc
	postRouter.Use(handlderProducts.MiddlewareProductValidation)

	deleteRouter := servMx.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.Delete)
	// executes middleware before it can go to the HandleFunc
	deleteRouter.Use(handlderProducts.MiddlewareProductValidation)

	//
	//   create web server
	server := &http.Server{
		Addr:         ":9090", // on port 9090
		Handler:      servMx,  // using defined serv mux
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	// start server on a separate process (non blocking)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	//signal.Notify(sigChan, os.Kill)
	// waiting on main process for termination signals to greacefully handle the event
	// potentially close any db connections etc
	sig := <-sigChan // blocking until receives a signal from sigChan
	log.Println("Received terminate, greaceful shotdown", sig)
	tcx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tcx)
}
