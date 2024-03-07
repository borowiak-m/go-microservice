package data

import (
	"fmt"
	"time"
)

// Product defines the structure for the API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products is a collection of Product
type Products []*Product

func GetProducts() Products {
	return productList
}

func GetProductByID(id int) (*Product, error) {
	index := findIndexByProductID(id)
	if index == -1 {
		return nil, ErrProductNotFound
	}
	return productList[index], nil
}

func AddProduct(prod *Product) {
	maxID := productList[len(productList)-1].ID
	fmt.Println("[DEBUG] maxID:", maxID)
	prod.ID = maxID + 1
	productList = append(productList, prod)
}

func UpdateProduct(prod *Product) error {
	fmt.Println("[DEBUG] looking for id:", prod.ID)
	index := findIndexByProductID(prod.ID)
	if index == -1 {
		return ErrProductNotFound
	}

	// update product
	productList[index] = prod
	return nil
}

func DeleteProduct(id int) error {
	index := findIndexByProductID(id)
	if index == -1 {
		return ErrProductNotFound
	}

	err := removeItemFromProducts(index)
	if err != nil {
		return err
	}
	return nil
}

func removeItemFromProducts(index int) error {
	listLen := len(productList)
	productList[index] = productList[listLen-1]
	productList = productList[:listLen-1]
	return nil
}

func findIndexByProductID(id int) int {
	for i, prod := range productList {
		if prod.ID == id {
			return i
		}
	}
	return -1 // if not found
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// Static data for time being as collection of Product
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "lat-cof-1",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Long double espresso",
		Description: "Tall espresso coffee shot without milk",
		Price:       1.99,
		SKU:         "exp-cof-1",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
