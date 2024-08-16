package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the logged-in user
// @Tags User
// @Produce  json
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Router /users/profile [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*models.Claims)
	if !ok || claims == nil {
		http.Error(w, "Unauthorized access or claims missing.", http.StatusUnauthorized)
		return
	}

	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found."+result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile of the logged-in user
// @Tags User
// @Accept  json
// @Produce  json
// @Param   user body models.User true "User"
// @Success 200 {string} string "Profile updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /users/profile [put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	var input models.User
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	user.Email = input.Email
	user.Name = input.Name
	user.Surname = input.Surname
	user.UpdatedAt = time.Now()

	if result := database.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update profile.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile update successfully."})
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update the password of the logged-in user
// @Tags User
// @Accept  json
// @Produce  json
// @Param   passwordData body models.PasswordUpdateRequest true "Password Update Request"
// @Success 200 {string} string "Password updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /users/profile/password [put]
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)

	var user models.User
	if result := database.DB.First(&user, claims.UserID); result.Error != nil {
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}

	var passwordData map[string]string
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	oldPass := passwordData["old_password"]
	newPass := passwordData["new_password"]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPass))
	if err != nil {
		http.Error(w, "Old password is incorrect.", http.StatusUnauthorized)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if result := database.DB.Save(&user); result.Error != nil {
		http.Error(w, "Failed to update password.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully."})
}

// DeleteUserHandler godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags User
// @Param   id path int true "User ID"
// @Success 204 {string} string "User deleted successfully"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/{id}/delete [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id.", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := database.DB.Delete(&user, userID); result.Error != nil {
		http.Error(w, "Failed to delete user.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CloseAccount godoc
// @Summary Close user account
// @Description Close the account of the logged-in user
// @Tags User
// @Success 204 {string} string "Account closed successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to close account"
// @Router /users/close-account [delete]
func CloseAccount(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(models.Claims)

	var user models.User
	if result := database.DB.Delete(&user, claims.UserID); result.Error != nil {
		http.Error(w, "Failed to close account.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
