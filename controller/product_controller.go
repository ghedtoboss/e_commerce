package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"time"
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
