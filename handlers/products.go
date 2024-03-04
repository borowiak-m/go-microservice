package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/gorilla/mux"
)

type Products struct {
	log *log.Logger
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log}
}

func (prods *Products) GetProducts(reqW http.ResponseWriter, req *http.Request) {
	// get products
	lp := data.GetProducts()
	// encode to JSON format
	err := lp.ToJSON(reqW)
	if err != nil {
		http.Error(reqW, "Unable to marshall products to json", http.StatusInternalServerError)
	}
}

func (prods *Products) AddProduct(reqW http.ResponseWriter, req *http.Request) {
	prods.log.Println("POST request response")
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	prods.log.Println("POST we have the prod:", prod)
	data.AddProduct(&prod)
}

func (prods *Products) UpdateProducts(reqW http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(reqW, "Unable to parse Product id", http.StatusBadRequest)
		return
	}

	prods.log.Println("Handle PUT Product id", id)
	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(reqW, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(reqW, "Product not found", http.StatusInternalServerError)
		return
	}

}

type KeyProduct struct {
}

// middleware takes next func of type http.Handler that can be chained
// docs: https://github.com/gorilla/mux?tab=readme-ov-file#middleware
func (prods *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(req.Body)
		if err != nil {
			prods.log.Println("Error deserializing Product", err)
			http.Error(respW, "Unable to parse from JSON request body to Product", http.StatusBadRequest)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			prods.log.Println("[ERROR] validating product", err)
			http.Error(
				respW,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		next.ServeHTTP(respW, req)
	})
}
