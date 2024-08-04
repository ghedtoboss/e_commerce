package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(models.Claims)

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
