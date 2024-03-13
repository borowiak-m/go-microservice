package handlers

import (
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

// POST request function to handle a creation of a new product
func (prods *Products) CreateProduct(reqW http.ResponseWriter, req *http.Request) {
	prods.log.Println("POST product request response")
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	prods.log.Println("[DEBUG] inserting prod:", prod)
	data.AddProduct(&prod)
}
