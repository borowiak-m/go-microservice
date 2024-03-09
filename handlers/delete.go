package handlers

import (
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

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
