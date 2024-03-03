package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

type Products struct {
	log *log.Logger
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log}
}

func (prods *Products) ServeHTTP(reqW http.ResponseWriter, req *http.Request) {
	// get products
	lp := data.GetProducts()
	// encode to JSON format
	dat, err := json.Marshal(lp)
	if err != nil {
		http.Error(reqW, "Unable to marshall products to json", http.StatusInternalServerError)
	}

	reqW.Write(dat)
}
