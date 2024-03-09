package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/gorilla/mux"
)

type Products struct {
	log *log.Logger
	val *data.Validation
}

func NewProducts(log *log.Logger, val *data.Validation) *Products {
	return &Products{log, val}
}

func getProductId(req *http.Request) int {
	// parse id from url
	vars := mux.Vars(req)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}

type KeyProduct struct {
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
