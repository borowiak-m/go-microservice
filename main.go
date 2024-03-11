package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/borowiak-m/go-microservice/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// new logger
	newlogger := log.New(os.Stdout, "product-api", log.LstdFlags)
	// new validation
	newValidation := data.NewValidation()

	// HANDLERS
	//   create handler for products with logger
	handlderProducts := handlers.NewProducts(newlogger, newValidation)
	//   create handler for users with logger
	handlerUsers := handlers.NewUsers(newlogger, newValidation)

	// SERVER
	//   new serve mux
	servMx := mux.NewRouter()
	//   register handlerProduct as server for "/products" pattern
	getProdRouter := servMx.Methods(http.MethodGet).Subrouter()
	getProdRouter.HandleFunc("/products", handlderProducts.GetProducts)
	getProdRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.GetSingleProduct)

	putProdRouter := servMx.Methods(http.MethodPut).Subrouter()
	putProdRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.UpdateSingleProduct)
	// executes middleware before it can go to the HandleFunc
	putProdRouter.Use(handlderProducts.MiddlewareProductValidation)

	postProdRouter := servMx.Methods(http.MethodPost).Subrouter()
	postProdRouter.HandleFunc("/products", handlderProducts.CreateProduct)
	// executes middleware before it can go to the HandleFunc
	postProdRouter.Use(handlderProducts.MiddlewareProductValidation)

	deleteProdRouter := servMx.Methods(http.MethodDelete).Subrouter()
	deleteProdRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.Delete)
	// executes middleware before it can go to the HandleFunc
	deleteProdRouter.Use(handlderProducts.MiddlewareProductValidation)

	//   register handlerProduct as server for "/products" pattern
	getUserRouter := servMx.Methods(http.MethodGet).Subrouter()
	getUserRouter.HandleFunc("/users", handlerUsers.Get200)
	//getUserRouter.HandleFunc("/products/{id:[0-9]+}", handlderProducts.GetSingleUser)

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
	log.Println("Received terminate, greaceful shutdown", sig)
	tcx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tcx)
}
