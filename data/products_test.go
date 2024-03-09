package data

import (
	"bytes"
	"testing"
)

func TestProduct_InvalidSKUReturnsErr(t *testing.T) {
	prod := Product{
		Name:  "nics",
		Price: 90,
		SKU:   "2AAA",
	}
	valid := NewValidation()
	err := valid.Validate(prod)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_MissingNameReturnsErr(t *testing.T) {
	prod := Product{
		Price: 1.22,
	}

	valid := NewValidation()
	err := valid.Validate(prod)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_IncorrectPriceReturnsErr(t *testing.T) {
	prod := Product{
		Name:  "Dsfdsf",
		Price: -1.22,
	}

	valid := NewValidation()
	err := valid.Validate(prod)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_ValidProductSuccess(t *testing.T) {
	prod := Product{
		Name:  "Namelyname",
		Price: 90,
		SKU:   "2AAA-GGG-DDD",
	}
	valid := NewValidation()
	err := valid.Validate(prod)

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_ProductToJSONSuccess(t *testing.T) {
	prods := []*Product{
		{Name: "abc"},
	}

	byts := bytes.NewBufferString("")
	err := ToJSON(prods, byts)
	if err != nil {
		t.Fatal(err)
	}
}
