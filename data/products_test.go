package data

import "testing"

func TestCheckValidation(t *testing.T) {
	prod := &Product{
		Name:  "nics",
		Price: 90,
		SKU:   "2AAA-DD-A",
	}
	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
