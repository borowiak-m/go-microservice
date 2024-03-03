package handlers

import (
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
	if req.Method == http.MethodGet {
		prods.getProducts(reqW, req)
		return
	}
	if req.Method == http.MethodPost {
		prods.addProduct(reqW, req)
		return
	}

	// catch all - for all other methods
	reqW.WriteHeader(http.StatusMethodNotAllowed)
}

func (prods *Products) getProducts(reqW http.ResponseWriter, req *http.Request) {
	// get products
	lp := data.GetProducts()
	// encode to JSON format
	err := lp.ToJSON(reqW)
	if err != nil {
		http.Error(reqW, "Unable to marshall products to json", http.StatusInternalServerError)
	}
}

func (prods *Products) addProduct(reqW http.ResponseWriter, req *http.Request) {
	log.Println("POST request response")
}
