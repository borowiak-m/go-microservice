package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

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
