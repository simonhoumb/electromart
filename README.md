# ElectroMart ⚡️🛒

## Description 📜
Electro mart is a web application that allows users to buy electronics.
The application is built as a part of the course IDATG2204 Datamodellering og databasesystemer (2024 VÅR) at NTNU.

## How to run the app 🤨

## Usage 🤓

### Endpoints

The service provides the following endpoints:

```plaintext
/api/v1/products/
/api/v1/categories/
/api/v1/brands/
```
Note: where the user needs to provide an ID, it can be found using the `View all` endpoint. or the `Search` endpoints.

#### Products
The initial endpoint focuses on the management products stored in the database.

##### Add new product to database

Manages the registration of new product and its details. This includes brand, category, quantity in stock, price name and description. Before adding the user has to make sure that both the category and brand already exists in the database.

###### Request

```http
Method: POST
Path: /api/v1/products/
Content type: application/json
```

Body example:

```json lines
{
    "name": "Supatest",
    "brandName": "Supabrand",
    "categoryName": "Supacategory",
    "description": "Supatest Supatest Supadescription",
    "qtyInStock": 10,
    "price": 999
}
```

###### Response

The response to the POST request on the endpoint stores the product on the server and returns the associated ID. 

* Content type: application/json

Body (exemplary code for registered configuration):

```json
{
    "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6"
}
```

##### View a specific product

Enables retrieval of a specific product stored in the database.
###### Request

The following shows a request for an individual configuration identified by its ID.

```text
Method: GET
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product.

Example request:

```http request
/api/v1/product/7acef38c-0d18-11ef-96c4-fa163ecc81b6
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json
{
    "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
    "name": "Supatest",
    "brandName": "Apple",
    "categoryName": "Smartphones",
    "description": "Supatest Supatest Supadescription",
    "qtyInStock": 10,
    "price": 999
}
```

##### View all registered products

Enables retrieval of all products.

###### Request

A `GET` request to the endpoint should return all products including their IDs.

```text
Method: GET
Path: /api/v1/products/
```

###### Response

* Content type: `application/json`

Body (exemplary code):

```json lines
[
    {
        "id": "7acef38c-0d18-11ef-96c4-fa163ecc81b6",
        "name": "Supatest",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Supatest Supatest Supadescription",
        "qtyInStock": 10,
        "price": 999
    },
    {
        "id": "ca868913-0d16-11ef-96c4-fa163ecc81b6",
        "name": "iPhone 13",
        "brandName": "Apple",
        "categoryName": "Smartphones",
        "description": "Latest iPhone model",
        "qtyInStock": 100,
        "price": 9999.99
    }
]
```

The response should return a collection of return all stored products.

##### Update a specific product
Enables the replacing of specific product.

###### Request

The following shows a request for an updated product identified by its ID.

```
Method: PUT
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product and must match in the provided url and body.

Example request: ```/api/v1/products/ca86bc9b-0d16-11ef-96c4-fa163ecc81b6```

Body (exemplary code):

```json lines
{
        "id": "ca86bc9b-0d16-11ef-96c4-fa163ecc81b6",
        "name": "Galaxy Book Pro",
        "brandName": "Samsung",
        "categoryName": "Laptops",
        "description": "Latest Samsung laptop model",
        "qtyInStock": 100,
        "price": 1199.99
}
```

###### Response

This is the response to the change request.

* Status code: `204 No Content` if the update was succsessful.
* Body: empty

##### Delete a specific registered dashboard configuration

Enabling the deletion of a specific product.

###### Request

The following shows a request for deletion of an individual product identified by its ID. This update should lead
to a deletion of the product from the server.

```text
Method: DELETE
Path: /api/v1/products/{id}
```

* `id` is the ID associated with the specific product.

Example request:

```http request
/api/v1/products/ca86bc9b-0d16-11ef-96c4-fa163ecc81b6
```

###### Response

This is the response to the delete request.

* Status code: `204 No Content`.
* Body: empty

---


#### Categories

#### Brands


## Contact 📧

- [Eskil Alstad](mailto:eskil.alstad@ntnu.no)
- [Erik Bjørnsen](mailto:erbj@stud.ntnu.no)
- [Simon Husås Houmb](mailto:simon.h.houmb@ntnu.no)
- [Maja Melby](mailto:maja.melby@ntnu.no)
- [Jon André Solberg](mailto:jon.a.h.solberg@ntnu.no)
