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

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for a product with the specified quantity
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "Bearer token"
// @Param   product_id path int true "Product ID"
// @Param   body body models.OrderItem true "Order Item"
// @Success 200 {string} string "Order created successfully"
// @Failure 400 {string} string "Invalid id" / "Invalid input" / "Not available in the required quantity"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Failed to create order" / "Failed to create order item" / "Failed to begin transaction" / "Failed to commit transaction"
// @Router /orders/{product_id} [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*models.Claims)
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

	var orderItem models.OrderItem
	err = json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}
	orderItem.ProductID = product.ID
	orderItem.CreatedAt = time.Now()
	orderItem.UpdatedAt = time.Now()
	orderItem.Price = product.Price
	orderItem.Total = orderItem.Price * float64(orderItem.Quantity)

	if product.Stock < orderItem.Quantity {
		http.Error(w, "Not available in the required quantity.", http.StatusBadRequest)
		return
	}

	var order models.Order
	order.UserID = claims.UserID
	order.TotalAmount = orderItem.Total
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Veritabanı işlemi (transaction) başlatma
	tx := database.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to begin transaction."+tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Siparişi oluşturma
	if result := tx.Create(&order); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Failed to create order."+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Sipariş öğesini oluşturma
	orderItem.OrderID = order.ID
	if result := tx.Create(&orderItem); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Failed to create order item."+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Transaction'ı commit etme (işlemi tamamla)
	if tx.Commit().Error != nil {
		http.Error(w, "Failed to commit transaction.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully."})
}
