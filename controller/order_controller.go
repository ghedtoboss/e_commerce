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

	tx := database.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to begin transaction."+tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result := tx.Create(&order); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Failed to create order."+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	orderItem.OrderID = order.ID
	if result := tx.Create(&orderItem); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Failed to create order item."+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if tx.Commit().Error != nil {
		http.Error(w, "Failed to commit transaction.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully."})
}

// UpdateOrderStatus godoc
// @Summary Update the status of an order
// @Description Updates the status of an order by its ID
// @Tags Orders
// @Accept  json
// @Produce  json
// @Param   order_id path int true "Order ID"
// @Param   body body object true "Order Status Update"
// @Success 200 {object} map[string]string "Order status updated successfully"
// @Failure 400 {string} string "Invalid id" / "Invalid input"
// @Failure 404 {string} string "Order not found"
// @Failure 500 {string} string "Failed to update order status"
// @Router /orders/{order_id}/status [put]
func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID, err := strconv.Atoi(params["order_id"])
	if err != nil {
		http.Error(w, "Invalid id.", http.StatusBadRequest)
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input.", http.StatusBadRequest)
		return
	}

	tx := database.DB.Begin()

	var order models.Order
	if result := tx.First(&order, orderID); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Order not found.", http.StatusNotFound)
		return
	}

	order.Status = input.Status

	if result := tx.Save(&order); result.Error != nil {
		tx.Rollback()
		http.Error(w, "Failed to update order status.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status update successfully."})

}
