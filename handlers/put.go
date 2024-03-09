package handlers

import (
	"net/http"
	"strconv"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/gorilla/mux"
)

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
