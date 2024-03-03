package handlers

import (
	"log"
	"net/http"
)

type Greeting struct {
	log *log.Logger
}

func NewGreeting(log *log.Logger) *Greeting {
	return &Greeting{log}
}

func (grt *Greeting) ServeHTTP(reqW http.ResponseWriter, req *http.Request) {
	reqW.Write([]byte("Greetings"))
}
