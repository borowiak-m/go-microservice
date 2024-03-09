package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

type ValidationError struct {
	validator.FieldError
}

// Produce an error description string with a few params from the ValidationError
func (vErr ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		vErr.Namespace(),
		vErr.Field(),
		vErr.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into string slice
func (vErrs ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range vErrs {
		errs = append(errs, err.Error())
	}

	return errs
}

type Validation struct {
	validate *validator.Validate
}

func NewValidation() *Validation {
	// create new validator
	validate := validator.New()
	// register the sku tag against validateSKU func
	validate.RegisterValidation("sku", validateSKU)
	// return struct with new validator
	return &Validation{validate}
}

// Using the validator package to use struct tags for data validation
// returns true or false if struct conforms to set out requirements
func (vln *Validation) Validate(strc interface{}) ValidationErrors {
	err := vln.validate.Struct(strc)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)

	var returnErrs []ValidationError
	for _, err := range validationErrors {
		// cast FieldError to ValidationError and append
		vErr := ValidationError{err}
		returnErrs = append(returnErrs, vErr)
	}
	return returnErrs
}

func validateSKU(field validator.FieldLevel) bool {

	reggieTheFinder := regexp.MustCompile(`[A-Za-z]+-[A-Za-z]+-[A-Za-z]+`)
	matches := reggieTheFinder.FindAllString(field.Field().String(), -1)
	fmt.Println(matches)
	return len(matches) == 1
}
