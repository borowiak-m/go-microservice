package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func RunServer(muxRouter *mux.Router) *http.Server {
	//   get port from env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	//   create web server
	server := &http.Server{
		Addr:         (":" + port), // on port 9090
		Handler:      muxRouter,    // using defined serv mux
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
	return server

}
