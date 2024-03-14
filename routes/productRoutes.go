package routes

import (
	"log"
	"net/http"

	"github.com/borowiak-m/go-microservice/data"
	"github.com/borowiak-m/go-microservice/handlers"
	"github.com/gorilla/mux"
)

func RegisterProductRoutes(newlogger *log.Logger, newValidation *data.Validation, muxRouter *mux.Router) {
	//   create handler for products with logger
	handlerProducts := handlers.NewProducts(newlogger, newValidation)
	// create subrouters per method
	getProdRouter := muxRouter.Methods(http.MethodGet).Subrouter()
	postProdRouter := muxRouter.Methods(http.MethodPost).Subrouter()
	putProdRouter := muxRouter.Methods(http.MethodPut).Subrouter()
	deleteProdRouter := muxRouter.Methods(http.MethodDelete).Subrouter()
	// GET /products
	getProdRouter.HandleFunc("/products", handlerProducts.GetProducts)
	// GET /products/{id}
	getProdRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.GetSingleProduct)
	// POST /products
	postProdRouter.HandleFunc("/products", handlerProducts.CreateProduct)
	// POST /products/ Middleware: executes before it can go to the HandleFunc
	postProdRouter.Use(handlerProducts.MiddlewareProductValidation)
	// PUT /products/{id}
	putProdRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.UpdateSingleProduct)
	// PUT /products/{id} Middleware: executes before it can go to the HandleFunc
	putProdRouter.Use(handlerProducts.MiddlewareProductValidation)
	// DELETE /products/{id}
	deleteProdRouter.HandleFunc("/products/{id:[0-9]+}", handlerProducts.Delete)
	// DELETE /products/{id} Middleware: executes before it can go to the HandleFunc
	deleteProdRouter.Use(handlerProducts.MiddlewareProductValidation)
}
