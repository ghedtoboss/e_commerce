package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
// @Router /shop [post]
func CreateShop(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

	var existingShop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&existingShop); result.Error == nil {
		http.Error(w, "You already have a shop.", http.StatusBadRequest)
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

// GetShop godoc
// @Summary Get shop information
// @Description Get information of a shop by shop ID
// @Tags Shop
// @Produce  json
// @Param   shop_id path int true "Shop ID"
// @Success 200 {object} models.Shop
// @Failure 400 {string} string "Invalid shop id"
// @Failure 404 {string} string "Shop not found"
// @Router /shop/{shop_id} [get]
func GetShop(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shopID, err := strconv.Atoi(params["shop_id"])
	if err != nil {
		http.Error(w, "Invalid shop id.", http.StatusBadRequest)
		return
	}

	var shop models.Shop
	if result := database.DB.First(&shop, shopID); result.Error != nil {
		http.Error(w, "Shop not found.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shop)
}

// GetMyShop godoc
// @Summary Get my shop information
// @Description Get information of the logged-in user's shop
// @Tags Shop
// @Produce  json
// @Success 200 {object} models.Shop
// @Failure 404 {string} string "Shop not found"
// @Failure 500 {string} string "Failed to retrieve shop"
// @Router /shop/my [get]
func GetMyShop(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*models.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized access or claims missing.", http.StatusUnauthorized)
		return
	}

	var shop models.Shop
	if result := database.DB.Where("owner_id = ?", claims.UserID).First(&shop); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Shop not found.", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve shop.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shop)
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
// @Router /shop [put]
func UpdateShop(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

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
