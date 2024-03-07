package data

import (
	"bytes"
	"testing"
)

func TestProduct_InvalidSKUReturnsErr(t *testing.T) {
	prod := &Product{
		Name:  "nics",
		Price: 90,
		SKU:   "2AAA",
	}
	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_MissingNameReturnsErr(t *testing.T) {
	prod := Product{
		Price: 1.22,
	}

	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_IncorrectPriceReturnsErr(t *testing.T) {
	prod := Product{
		Name:  "Dsfdsf",
		Price: -1.22,
	}

	err := prod.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

func TestProduct_ValidProductSuccess(t *testing.T) {
	prod := &Product{
		Name:  "Namelyname",
		Price: 90,
		SKU:   "2AAA-GGG-DDD",
	}
	err := prod.Validate()

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
