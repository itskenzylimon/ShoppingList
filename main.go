package main

import (
	"example/ShoppingList/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseModel represents a response model.
// [StatusCode] is the status code of the response. i.e. 200, 400, 500
// [Success] is the status of the response, either true or false
// [Message] is the message of the response.
// [Data] is the data of the response.
type ResponseModel struct {
	StatusCode int    `json:"statusCode"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

// ValidateProduct represents a product in the shopping list.
// [Name] is the name of the product.
// [Description] is the description of the product.
// [Image] is the image of the product, a base64 string.
// [Price] is the price of the product.
type ValidateProduct struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}

// ValidateShopping represents a shopping list.
// [Name] is the name of the shopping list.
// [Description] is the description of the shopping list.
type ValidateShopping struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func main() {

	// Custom config
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "shoppingList",
	})

	models.SetUpDatabase()

	// Create a new response model, with default values
	responseModel := ResponseModel{
		Success:    true,
		StatusCode: 200,
		Message:    "",
		Data:       nil,
	}

	// Create a new product Route
	app.Post("/api/product", func(c *fiber.Ctx) error {

		// Declare a new validateProduct struct
		var validateProduct ValidateProduct

		// Parse body into validateProduct
		if err := c.BodyParser(&validateProduct); err != nil {
			responseModel.Data = err
			responseModel.Message = "Failed to parse JSON body"
			responseModel.StatusCode = 422
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// validate request body using validator
		error := validator.New().Struct(validateProduct); if error != nil {
			responseModel.Data = error
			responseModel.Message = "Invalid parameter"
			responseModel.StatusCode = 422
			responseModel.Success = false
			return c.Status(fiber.StatusUnauthorized).JSON(responseModel)
		}

		// Create product
		product := models.CreateProduct(
			validateProduct.Name,
			validateProduct.Description,
			validateProduct.Image,
			validateProduct.Price)

		responseModel.Data = product
		responseModel.Message = "Product created successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Read a single product or get products Route
	app.Get("/api/product", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Query("id", "0"))
		// Find product by id
		products := models.GetProducts(id)

		responseModel.Data = products
		responseModel.Message = "Product retrieved successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Update a product
	app.Put("/api/product/:id", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Params("id", "0"))
		if id == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Product Id"
			responseModel.StatusCode = 422
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// /// Get product by id
		// var product []models.Product
		product := models.GetProducts(id)
		if product == nil {
			responseModel.Data = nil
			responseModel.Message = "Product not found"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		newProduct := product[0]
		newProduct.Name = c.Query("name", product[0].Name)
		newProduct.Description = c.Query("description", product[0].Description)
		newProduct.Price, _ = strconv.Atoi(c.Query("price", strconv.Itoa(product[0].Price)))

		// Update product
		models.UpdateProduct(newProduct.ID, newProduct)

		responseModel.Data = newProduct
		responseModel.Message = "Product updated successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Read a product or get products Route
	app.Delete("/api/product/:id", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Params("id", "0"))
		if id == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Product Id"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		/// Get product by id
		product := models.GetProducts(id)
		if product == nil {
			responseModel.Data = nil
			responseModel.Message = "Product not found"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// Delete product
		models.DeleteProduct(product[0].ID)

		responseModel.Data = nil
		responseModel.Message = "Product deleted successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Create a new shopping name Route
	app.Post("/api/shoppingList", func(c *fiber.Ctx) error {

		// Declare a new Product struct.
		var validateShopping ValidateShopping

		// Parse body into product struct
		if err := c.BodyParser(&validateShopping); err != nil {
			responseModel.Data = err
			responseModel.Message = "Failed to parse JSON body"
			responseModel.StatusCode = 422
			responseModel.Success = false
			// Return status 400 bad request
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// validate request body using validator
		error := validator.New().Struct(validateShopping)
		if error != nil {
			responseModel.Data = error
			responseModel.Message = "Invalid parameter"
			responseModel.StatusCode = 422
			responseModel.Success = false
			return c.Status(fiber.StatusUnauthorized).JSON(responseModel)
		}

		// Create shoppingList
		shoppingList := models.CreateShoppingList(
			validateShopping.Name,
			validateShopping.Description)

		responseModel.Data = shoppingList
		responseModel.Message = "Added Shopping List successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Read shopping list or get shopping Route
	app.Get("/api/fetchShoppingList", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Query("id", "0"))
		// Find product by id
		shoppingList, err := models.FetchShoppingList(id)
		if err != nil {
			responseModel.Data = err
			responseModel.Message = "Failed get shopping list"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		responseModel.Data = shoppingList
		responseModel.Message = "Shopping list Loaded successfully"
		responseModel.StatusCode = 200

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Update a shopping
	app.Put("/api/shoppingList/:id", func(c *fiber.Ctx) error {

		// Get Shopping List Id from URL params and convert string to int
		id, _ := strconv.Atoi(c.Params("id", "0"))
		if id == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Shopping List Id"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		/// Get Shopping List by id
		shoppingList := models.GetShoppingList(id)
		if shoppingList == nil {
			responseModel.Data = nil
			responseModel.Message = "Shopping List not found"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		newShoppingList := shoppingList[0]
		newShoppingList.Name = c.Query("name", newShoppingList.Name)
		newShoppingList.Description = c.Query("description", newShoppingList.Description)

		// Update Shopping List
		models.UpdateShoppingList(newShoppingList.ID, newShoppingList)

		responseModel.Data = newShoppingList
		responseModel.Message = "Successfully updated shopping list"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Read a shopping or get shopping Route
	app.Delete("/api/shoppingList/:id", func(c *fiber.Ctx) error {

		// Get Shopping List Id from URL params and convert string to int
		id, _ := strconv.Atoi(c.Params("id", "0"))
		if id == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Product Id"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		/// Get Shopping List by id
		product := models.GetProducts(id)
		if product == nil {
			responseModel.Data = nil
			responseModel.Message = "Product not found"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// Delete Shopping List
		models.DeleteShoppingList(product[0].ID)

		responseModel.Data = nil
		responseModel.Message = "Shopping List deleted successfully"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	// Add product to shopping list
	app.Put("/api/shoppingList/:id/product/:productID/attach", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Params("id", "0"))
		productID, _ := strconv.Atoi(c.Params("productID", "0"))

		// Check if id and productID are valid
		if id == 0 || productID == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Shopping List Id or Product Id"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// Add product to shopping list
		models.AddProductToShoppingList(id, productID)

		responseModel.Data = nil
		responseModel.Message = "Successfully added to shopping list"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	//Remove product from shopping list
	app.Put("/api/shoppingList/:id/product/:productID/remove", func(c *fiber.Ctx) error {

		// Get Product Id from URL params and convert string to uint
		id, _ := strconv.Atoi(c.Params("id", "0"))
		productID, _ := strconv.Atoi(c.Params("productID", "0"))

		fmt.Println(id, productID)

		// Check if id and productID are valid
		if id == 0 || productID == 0 {
			responseModel.Data = nil
			responseModel.Message = "Invalid Shopping List or Product"
			responseModel.StatusCode = 404
			responseModel.Success = false
			return c.Status(responseModel.StatusCode).JSON(responseModel)
		}

		// Remove product to shopping list
		models.RemoveProductToShoppingList(id, productID)

		responseModel.Data = nil
		responseModel.Message = "Successfully Removed to shopping list"
		responseModel.StatusCode = 200
		responseModel.Success = true

		return c.Status(responseModel.StatusCode).JSON(responseModel)
	})

	app.Listen(":3000")

}
