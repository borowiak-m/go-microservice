package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/borowiak-m/go-microservice/handlers"
)

func main() {
	// new logger
	newlogger := log.New(os.Stdout, "product-api", log.LstdFlags)

	// HANDLERS
	//   create handler for root with new logger
	handlerRoot := handlers.NewRoot(newlogger)
	//   create handler for greetings with logger
	handlerGreeting := handlers.NewGreeting(newlogger)
	//   create handler for products with logger
	handlderProducts := handlers.NewProducts(newlogger)

	// SERVER
	//   new serve mux
	servMx := http.NewServeMux()
	//   register handlerRoot as server for "/" pattern
	servMx.Handle("/", handlerRoot)
	//   register handlerGreeting as server for "/greeting" pattern
	servMx.Handle("/greeting", handlerGreeting)
	//   register handlerProduct as server for "/products" pattern
	servMx.Handle("/products", handlderProducts)
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
