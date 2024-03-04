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

**Success response:**

- Status code: ``` 200 OK```, Response Body: A JSON array of product objects, each containing ID, Name, Description, Price, and SKU.

Example:

```
[
  {
    "id": 1,
    "name": "Latte",
    "description": "Frothy milky coffee",
    "price": 2.45,
    "sku": "lat-cof-1"
  },
  {
    "id": 2,
    "name": "Espresso",
    "description": "Strong coffee without milk",
    "price": 1.99,
    "sku": "esp-cof-2"
  }
]

```

**Error responses:**

- Status code: ```500 Internal Server Error```, Error message: ```"Unable to marshall products to JSON"```

### Add Product

POST /products

**Success response:**

- Status code: ``` 200 OK```, Response Body: No body returned

**Error Responses:**

- Status Code: ```400 Bad Request```, Error Message: "Unable to parse from JSON request body to Product." or "Error validating product: [validation error message]."

### Update Product 

PUT /products/{id}

Updates the details of an existing product by ID. 
Requires a JSON body with the fields you wish to update.

**Success response:**

- Status code: ``` 200 OK```, Response Body: No body returned

**Error responses:**

- Status Code: ```400 Bad Request```, Error Message: "Unable to parse Product id." or "Unable to parse from JSON request body to Product." or "Error validating product: [validation error message]."
- Status Code: ```404 Not Found```, Error Message: "Product not found."
- Status Code: ```500 Internal Server Error```, Error Message: "Product update failed."

### Examples

Adding a new product:

POST /products

  Content-Type: application/json
```
{
  "name": "Espresso",
  "description": "Strong coffee without milk",
  "price": 1.99,
  "sku": "esp-cof-yyt"
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
  "sku": "esp-cof-yyy"
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

## SKU Validation and Business Logic Considerations

While our API enforces a specific format for the SKU to ensure consistency and readability, it's important to note that we do not validate the uniqueness of the SKU within our system. 
This decision is intentionally aligned with the flexibility needed to accommodate various business requirements, especially in scenarios where product grouping plays a crucial role.

### Business Requirements and Product Grouping

In certain industries, such as fresh produce distribution, the differentiation between products might not solely rely on the SKU. For example, a range of tomatoes could be differentiated based on a combination of factors like SKU, weight, or case size. 
In these cases, products from different vendors may be treated as a single SKU, despite variations in pallet or container size.

This approach allows our API to be adaptable to diverse business models, where the level of product grouping and differentiation criteria may vary significantly. Here are a few key points to consider:

- Flexibility in Product Cataloging: By not enforcing SKU uniqueness, our API provides businesses the flexibility to catalog their products based on their unique operational and logistical requirements.

- Adaptability to Business Models: This approach supports various business models and distribution strategies, particularly in sectors where products are similar but differentiated by secondary attributes.

- Custom Implementation for Uniqueness: Should your business logic require unique SKUs for each product variant, this constraint would need to be implemented at the application level, taking into consideration the specific needs and rules of your operational model.









