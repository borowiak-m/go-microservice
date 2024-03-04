package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
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

// Using the validator package to use struct tags for data validation
// returns true or false if struct conforms to set out requirements
func (prod *Product) Validate() error {
	validate := validator.New()
	// register a validation func that will run on "sku" tag
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(prod)
}

func validateSKU(field validator.FieldLevel) bool {

	reggieTheFinder := regexp.MustCompile(`[A-Za-z]+-[A-Za-z]+-[A-Za-z]+`)
	matches := reggieTheFinder.FindAllString(field.Field().String(), -1)
	fmt.Println(matches)
	return len(matches) == 1
}

// Products is a collection of Product
type Products []*Product

// Serializes contents of the collection to JSON
// NewEncoder has better performance than json.Unmarshall
// due to no usage of interim buffer in memory
func (prods *Products) ToJSON(wrt io.Writer) error {
	encoder := json.NewEncoder(wrt)
	return encoder.Encode(prods)
}

// Decoding struct from JSON
func (prod *Product) FromJSON(re io.Reader) error {
	decoder := json.NewDecoder(re)
	return decoder.Decode(prod)
}

func GetProducts() Products {
	return productList
}

func AddProduct(prod *Product) {
	prod.ID = getNextId()
	productList = append(productList, prod)
}

func UpdateProduct(id int, prod *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}

	prod.ID = id
	productList[pos] = prod

	return nil
}

var ErrProductNotFound = fmt.Errorf("Product not found")

// Find product by id, used in PUT request func to return a single item
// and update its params from request
func findProduct(id int) (*Product, int, error) {
	for index, prod := range productList {
		if prod.ID == id {
			return prod, index, nil
		}
	}

	return &Product{}, -1, ErrProductNotFound
}

func getNextId() int {
	lastItem := productList[len(productList)-1]
	return lastItem.ID + 1
}

// Static data for time being as collection of Product
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "lat-cof-1",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Long double espresso",
		Description: "Tall espresso coffee shot without milk",
		Price:       1.99,
		SKU:         "exp-cof-1",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
