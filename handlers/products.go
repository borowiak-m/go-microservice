package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/borowiak-m/go-microservice/data"
)

type Products struct {
	log *log.Logger
}

func NewProducts(log *log.Logger) *Products {
	return &Products{log}
}

func (prods *Products) ServeHTTP(reqW http.ResponseWriter, req *http.Request) {
	prods.log.Println("ServeHTTP response", req.URL.Path)
	if req.Method == http.MethodGet {
		prods.getProducts(reqW, req)
		return
	}
	if req.Method == http.MethodPost {
		prods.addProduct(reqW, req)
		return
	}

	if req.Method == http.MethodPut {
		prods.log.Println("PUT response", req.URL.Path)
		reg := regexp.MustCompile(`/([0-9]+)`)
		grp := reg.FindAllStringSubmatch(req.URL.Path, -1)
		prods.log.Println(grp)
		if len(grp) != 1 {
			prods.log.Println("Invalid URI more than one id")
			http.Error(reqW, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(grp[0]) != 2 {
			prods.log.Println("Invalid URI more than one capture group")
			http.Error(reqW, "Invalid URI", http.StatusBadRequest)
			return
		}

		idStr := grp[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			prods.log.Println("Invalid URI unable to convert to number", idStr)
			http.Error(reqW, "Couldn't convert id to int", http.StatusBadRequest)
			return
		}
		prods.log.Println("Got id", id)
		prods.updateProducts(id, reqW, req)

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
	prods.log.Println("POST request response")
	// define an empty Product
	prod := &data.Product{}
	// get params for new product from request body
	err := prod.FromJSON(req.Body)
	// Be a good citizen and check error
	if err != nil {
		http.Error(reqW, "Unable to parse from JSON request body to Product", http.StatusBadRequest)
	}
	// submit new Product to data persist
	prods.log.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (prods *Products) updateProducts(id int, reqW http.ResponseWriter, req *http.Request) {
	prods.log.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJSON(req.Body)
	if err != nil {
		http.Error(reqW, "Unable to parse from JSON request body to Product", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(reqW, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(reqW, "Product not found", http.StatusInternalServerError)
		return
	}

}
