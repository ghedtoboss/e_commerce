package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// CreateShop godoc
// @Summary Create a new shop
// @Description Create a new shop for the logged-in seller
// @Tags Shop
// @Accept  json
// @Produce  json
// @Param   shop body models.Shop true "Shop"
// @Success 201 {string} string "Shop created successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to create shop"
// @Router /shops [post]
func CreateShop(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(models.Claims)

	var existingShop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&existingShop); result.Error == nil {
		http.Error(w, "You already have a shop.", http.StatusBadRequest)
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		http.Error(w, "Failed to check existing shop.", http.StatusInternalServerError)
		return
	}

	var shop models.Shop
	err := json.NewDecoder(r.Body).Decode(&shop)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	shop.OwnerID = claims.UserID
	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()

	if result := database.DB.Create(&shop); result.Error != nil {
		http.Error(w, "Failed to create shop.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Shop created successfully."})
}

// UpdateShop godoc
// @Summary Update shop information
// @Description Update the information of the logged-in user's shop
// @Tags Shop
// @Accept  json
// @Produce  json
// @Param   shop body models.Shop true "Shop"
// @Success 200 {string} string "Shop updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Shop not found"
// @Failure 500 {string} string "Failed to update shop"
// @Router /shops [put]
func UpdateShop(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(models.Claims)

	var shop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&shop); result.Error != nil {
		http.Error(w, "Shop not found.", http.StatusNotFound)
		return
	}

	var input models.Shop
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	shop.Name = input.Name
	shop.UpdatedAt = time.Now()

	if result := database.DB.Save(&shop); result.Error != nil {
		http.Error(w, "Failed to update shop.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Shop updated successfully."})
}
