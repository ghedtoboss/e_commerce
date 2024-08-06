package controller

import (
	"bytes"
	"e_commerce/database"
	"e_commerce/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	t.Run("Successfull Registiration", func(t *testing.T) {
		user := models.User{
			Name:     "John",
			Surname:  "Doe",
			Email:    "john.doe@example.com",
			Password: "password",
			Role:     "customer",
		}
		userJSON, _ := json.Marshal(user)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RegisterHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v wantr %v", status, http.StatusCreated)
		}

		var createdUser models.User
		err = json.NewDecoder(rr.Body).Decode(&createdUser)
		if err != nil {
			t.Fatal(err)
		}

		if createdUser.Email != user.Email {
			t.Errorf("expected user email to be %v, got %v", user.Email, createdUser.Email)
		}
	})

	t.Run("Invalid input", func(t *testing.T) {
		invalidUserJSON := []byte(`{"name":"", "email":"invalid"}`)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(invalidUserJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RegisterHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Missing Role", func(t *testing.T) {
		user := models.User{
			Name:     "Jane",
			Surname:  "Doe",
			Email:    "jane.doe@example.com",
			Password: "password",
		}

		userJSON, _ := json.Marshal(user)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(RegisterHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}
