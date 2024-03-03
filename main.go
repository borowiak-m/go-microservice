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

	// HANDLERS
	//   create handler for root with new logger
	handlerRoot := handlers.NewRoot(newlogger)
	//   create handler for greetings with logger
	handlerGreeting := handlers.NewGreeting(newlogger)

	// SERVER
	//   new serve mux
	servMx := http.NewServeMux()
	//   register handlerRoot as server for "/" pattern
	servMx.Handle("/", handlerRoot)
	//   register handlerGreeting as server for "/greeting" pattern
	servMx.Handle("/greeting", handlerGreeting)
	//   create web server on port 9090 using defined serv mux
	http.ListenAndServe(":9090", servMx)
}
