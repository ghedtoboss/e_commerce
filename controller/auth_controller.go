package controller

import (
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with username, password, email, and role
// @Tags User
// @Accept  json
// @Produce  json
// @Param   user body models.User true "User"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /users/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	if user.Role == "" {
		http.Error(w, "Role is Required.", http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPass)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if result := database.DB.Create(&user); result.Error != nil {
		http.Error(w, "Failed to create user: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

}

// LoginHandler godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags User
// @Accept  json
// @Produce  json
// @Param   login body models.LoginRequest true "Login Request"
// @Success 200 {string} string "Logged in successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid email or password"
// @Failure 500 {string} string "Internal server error"
// @Router /users/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var reqUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := database.DB.Where("email = ?", reqUser.Email).First(&user); result.Error != nil {
		http.Error(w, "Invalid email or password.", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		http.Error(w, "Invalid email or password.", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(168 * time.Hour)
	claims := &models.Claims{
		Email:  user.Email,
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
