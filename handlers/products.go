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
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log}
}

func (prods *Products) GetProducts(reqW http.ResponseWriter, req *http.Request) {
	// get products
	allProds := data.GetProducts()
	// encode to JSON format
	err := data.ToJSON(allProds, reqW)
	if err != nil {
		http.Error(reqW, "Unable to marshall products to json", http.StatusInternalServerError)
	}
}

func (prods *Products) GetSingleProduct(respW http.ResponseWriter, req *http.Request) {
	id := getProductId(req)

	prods.log.Println("[DEBUG] get record id:", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		prods.log.Println("[ERROR] fetching product:", err)

		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	default:
		prods.log.Println("[ERROR] fetching product:", err)

		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}

	if err = data.ToJSON(prod, respW); err != nil {
		prods.log.Println("[ERROR] serializing product", err)
	}

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

// POST request function to handle a creation of a new product
func (prods *Products) CreateProduct(reqW http.ResponseWriter, req *http.Request) {
	prods.log.Println("POST request response")
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	prods.log.Println("[DEBUG] inserting prod:", prod)
	data.AddProduct(&prod)
}

// DELETE /products/{id}
func (prods *Products) Delete(respW http.ResponseWriter, req *http.Request) {
	// get product
	id := getProductId(req)
	// log
	prods.log.Println("[DEBUG] deleting record id", id)
	// process deletion
	err := data.DeleteProduct(id)
	// error handle
	// if prod not found
	if err == data.ErrProductNotFound {
		prods.log.Println("[ERROR] deleting product id doesn't exist")

		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}
	// if other error
	if err != nil {
		prods.log.Println("[ERROR] deleting product")

		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, respW)
		return
	}

	respW.WriteHeader(http.StatusNoContent)
}

// PUT request function to handle updating items parameters fetched by id variable
// via context storage with Gorilla framework
func (prods *Products) UpdateSingleProduct(respW http.ResponseWriter, req *http.Request) {
	prods.log.Println("[DEBUG] UpdateSingleProduct starting")
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(respW, "[ERROR] Unable to parse Product id", http.StatusBadRequest)
		return
	}

	prods.log.Println("[DEBUG] Handle PUT Product id", id)
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	// overwrite prod id in case body prod id != URL prod id
	prod.ID = id
	prods.log.Println("New product data", prod)
	err = data.UpdateProduct(&prod)
	// if product not found
	if err == data.ErrProductNotFound {
		prods.log.Println("[ERROR] product not found", err)

		respW.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, respW)
		return
	}
	// if other error when updating
	if err != nil {
		prods.log.Println("[ERROR] updating product", err)

		respW.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: "Product update FAILED"}, respW)
		return
	}
	// success
	respW.WriteHeader(http.StatusNoContent)

}

type KeyProduct struct {
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}
