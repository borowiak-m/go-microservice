package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
)

// middleware takes next func of type http.Handler that can be chained
// docs: https://github.com/gorilla/mux?tab=readme-ov-file#middleware
func (prods *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		prod := data.Product{}

		err := data.FromJSON(&prod, req.Body)
		if err != nil {
			prods.log.Println("[ERROR] deserializing Product", err)
			respW.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, respW)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			prods.log.Println("[ERROR] validating product", err)
			respW.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&GenericError{Message: fmt.Sprintf("Error validating product: %s", err)}, respW)
			return
		}

		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		next.ServeHTTP(respW, req)
	})
}
