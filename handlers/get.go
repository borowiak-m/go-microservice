package handlers

import (
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

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
