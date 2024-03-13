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
	handlerProducts := handlers.NewProducts(newlogger, newValidation)
	//   create handler for users with logger
	handlerUsers := handlers.NewUsers(newlogger, newValidation)

	// SERVER
	//   new serve mux
	servMx := mux.NewRouter()
	//   register handlerProduct as server for "/products" pattern
	getRouter := servMx.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", handlerProducts.GetProducts)
	getRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.GetSingleProduct)
	getRouter.HandleFunc("/users", handlerUsers.Get200)
	getRouter.HandleFunc("/users/{id:[0-9]+}", handlerUsers.GetSingleUser)
	getRouter.Use(handlerUsers.MiddlewareUserAuth)

	putRouter := servMx.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.UpdateSingleProduct)
	// executes middleware before it can go to the HandleFunc
	putRouter.Use(handlerProducts.MiddlewareProductValidation)

	postProdRouter := servMx.Methods(http.MethodPost).Subrouter()
	postProdRouter.HandleFunc("/products", handlerProducts.CreateProduct)
	// executes middleware before it can go to the HandleFunc
	postProdRouter.Use(handlerProducts.MiddlewareProductValidation)

	postUserRouter := servMx.Methods(http.MethodPost).Subrouter()
	postUserRouter.HandleFunc("/users/signup", handlerUsers.CreateUser)
	postUserRouter.HandleFunc("/users/login", handlerUsers.Get200)
	// executes middleware before it can go to the HandleFunc
	postUserRouter.Use(handlerUsers.MiddlewareUserAuth)
	postUserRouter.Use(handlerUsers.MiddlewareUserValidation)

	deleteRouter := servMx.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.Delete)
	// executes middleware before it can go to the HandleFunc
	deleteRouter.Use(handlerProducts.MiddlewareProductValidation)

	//   get port from env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	//   create web server
	server := &http.Server{
		Addr:         (":" + port), // on port 9090
		Handler:      servMx,       // using defined serv mux
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
