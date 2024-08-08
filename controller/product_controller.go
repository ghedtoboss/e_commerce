package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// AddProduct godoc
// @Summary Add a new product
// @Description Add a new product to the shop of the logged-in user
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product details"
// @Success 201 {object} models.Product
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Shop not found"
// @Failure 500 {string} string "Failed to create product"
// @Router /products [post]
func AddProduct(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var shop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&shop); result.Error != nil {
		http.Error(w, "Shop not found.", http.StatusNotFound)
		return
	}

	product.ShopID = shop.ID
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if result := database.DB.Create(&product); result.Error != nil {
		http.Error(w, "Failed to create product.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update the details of an existing product
// @Tags Products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param product body models.Product true "Updated product details"
// @Success 200 {string} string "Product updated successfully."
// @Failure 400 {string} string "Invalid id or input"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Failed to update product"
// @Router /products/{product_id} [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["product_id"])
	if err != nil {
		http.Error(w, "Invalid id.", http.StatusBadRequest)
		return
	}

	var input models.Product
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		http.Error(w, "Product not found.", http.StatusNotFound)
		return
	}

	product.Stock = input.Stock
	product.Price = input.Price
	product.UpdatedAt = time.Now()
	product.ImageUrl = input.ImageUrl
	product.Description = input.Description
	product.Name = input.Name

	if result := database.DB.Save(&product); result.Error != nil {
		http.Error(w, "Failed to update product.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated successfully."})

}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get the details of a product by its ID
// @Tags Products
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {string} string "Invalid id"
// @Failure 404 {string} string "Product not found"
// @Router /products/{product_id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["product_id"])
	if err != nil {
		http.Error(w, "Invalid id.", http.StatusBadRequest)
		return
	}

	var product models.Product
	if result := database.DB.First(&product, productID); result.Error != nil {
		http.Error(w, "Product not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// GetProducts godoc
// @Summary Get all products
// @Description Get the details of all products
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 404 {string} string "Products not found"
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product
	if result := database.DB.Find(&products); result.Error != nil {
		http.Error(w, "Products not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// GetProductsByShop godoc
// @Summary Get products by shop ID
// @Description Get the details of all products for a specific shop
// @Tags Products
// @Produce json
// @Param shop_id path int true "Shop ID"
// @Success 200 {array} models.Product
// @Failure 400 {string} string "Invalid id"
// @Failure 404 {string} string "Products not
func GetProductsByShop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shopID, err := strconv.Atoi(params["shop_id"])
	if err != nil {
		http.Error(w, "Invalid id.", http.StatusBadRequest)
		return
	}

	var products []models.Product
	if result := database.DB.Where("shop_id = ?", shopID).Find(&products); result.Error != nil {
		http.Error(w, "Products not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// GetProductsByMyShop godoc
// @Summary Get products by the logged-in user's shop
// @Description Get the details of all products for the logged-in user's shop
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 404 {string} string "Shop or products not found"
// @Router /myshop/products [get]
func GetProductsByMyShop(w http.ResponseWriter, r *http.Request) {
	log.Println("fonksiyon çalıştı.")
	claims := r.Context().Value("user").(*models.Claims)

	var shop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&shop); result.Error != nil {
		http.Error(w, "Shop not found.", http.StatusInternalServerError)
		return
	}

	var products []models.Product
	if result := database.DB.Where("shop_id = ?", shop.ID).Find(&products); result.Error != nil {
		http.Error(w, "Products not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
