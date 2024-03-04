# Product Maintenance API
## Introduction

This project is a Golang-based RESTful API designed to manage product data. It allows for the creation, retrieval, updating, and validation of product information. Utilizing the powerful Gorilla Mux router and the go-playground validator, it provides a robust and efficient way to handle product data with validated input.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.15 or later)
- Gorilla Mux (documentation: https://pkg.go.dev/github.com/gorilla/mux)
- Go Playground Validator (documentation: https://pkg.go.dev/github.com/go-playground/validator/v10)

## Installation

Clone the repository to your local machine: 

``` git clone https://github.com/borowiak-m/gorilla-products-maintenance```

Navigate to the project directory and install dependencies:

```
go get -u github.com/gorilla/mux
go get -u github.com/go-playground/validator/v10
```
To start the server, run: 

``` go run main.go ```

The server will start on port 9090. You can access the API at http://localhost:9090/

## API Endpoints

### Get Products

GET /products

### Add Product

POST /products

### Update Product 

PUT /products/{id}

Updates the details of an existing product by ID. 
Requires a JSON body with the fields you wish to update.

### Examples

Adding a new product:

POST /products
Content-Type: application/json
```
{
  "name": "Espresso",
  "description": "Strong coffee without milk",
  "price": 1.99,
  "sku": "esp-cof-2"
}
```

Updating a product:

PUT /products/1
Content-Type: application/json
```
{
  "name": "Double Espresso",
  "description": "Extra strong espresso coffee",
  "price": 2.49,
  "sku": "esp-cof-2"
}
```

## Data Model and Validation

### Product structure
Definition of the product structure:
- ID (int): A unique identifier for each product. Automatically managed by the system.
- Name (string): The name of the product. This field is required for every product entry.
- Description (string): A detailed description of the product.
- Price (float32): The price of the product in USD. The value must be greater than 0.
- SKU (string): The Stock Keeping Unit, a unique code to identify each product variant, adhering to a specific format.
- CreatedOn, UpdatedOn, DeletedOn (string): Timestamps for tracking the lifecycle of the product. These fields are for internal use and not exposed in the API responses.

Response json example:
```
{
    "id": 1,
    "name": "Latte",
    "description": "Frothy milky coffee",
    "price": 2.45,
    "sku": "lat-cof-abc"
  }
```

### Validation method

Our API utilizes the validator package to ensure that product data meets our predefined criteria before it's processed or stored. 
This validation process enhances data integrity and prevents common data entry errors. Mind you, in this current iteration of the project, the data persistency happens only in memory but validation steps used apply to any data storage solution.

- Required Fields: Certain fields like Name and SKU are marked as required. The API rejects any product creation or update requests missing these fields.
- Price Validation: The Price field must be greater than 0. This is crucial for ensuring that all products have valid pricing information.
- Custom SKU Validation: The SKU field must match a specific pattern, validated through a custom function.

### Custom SKU validation

The SKU field uses a custom validator to ensure that each product's SKU adheres to a specific format: 
three groups of letters separated by hyphens (e.g., abc-def-ghi). 
This format ensures consistency and avoids conflicts in product identification. 
Here's an overview of the validation process:
- A regular expression ([A-Za-z]+-[A-Za-z]+-[A-Za-z]+) is used to match the SKU format.
- The custom validator, validateSKU, checks if the SKU field conforms to this pattern. Only SKUs matching the pattern pass the validation.
- This validation is part of the product's Validate method, ensuring that no product is added or updated in the system without a properly formatted SKU.






