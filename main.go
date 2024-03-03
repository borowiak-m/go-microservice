package main

import (
	"log"
	"net/http"
	"os"

	"github.com/borowiak-m/go-microservice/handlers"
)

func main() {
	// new logger
	newlogger := log.New(os.Stdout, "product-api", log.LstdFlags)
	// create handler for root with new logger
	handlerRoot := handlers.NewRoot(newlogger)
	// new serve mux
	servMx := http.NewServeMux()
	// register handler as server for "/" pattern
	servMx.Handle("/", handlerRoot)
	// create web server on port 9090 using defined serv mux
	http.ListenAndServe(":9090", servMx)
}
