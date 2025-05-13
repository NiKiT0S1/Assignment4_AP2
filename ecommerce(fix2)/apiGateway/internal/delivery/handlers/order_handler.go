package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	// You'll need to create these proto imports
	"apiGateway/internal/proto/order"
)

// OrderClient interface for the order service
type OrderClient interface {
	CreateOrder(ctx context.Context, order *order.Order, opts ...grpc.CallOption) (*order.Order, error)
	GetOrder(ctx context.Context, id *order.OrderID, opts ...grpc.CallOption) (*order.Order, error)
	UpdateOrderStatus(ctx context.Context, order *order.Order, opts ...grpc.CallOption) (*order.Order, error)
	ListOrdersByUser(ctx context.Context, req *order.ListOrdersRequest, opts ...grpc.CallOption) (*order.OrderList, error)
}

// OrderHandler handles HTTP requests for order service
type OrderHandler struct {
	client OrderClient
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(client OrderClient) *OrderHandler {
	return &OrderHandler{client: client}
}

// GetOrders returns all orders for the current user
func (h *OrderHandler) GetOrders(c *gin.Context) {
	// Get user ID from the context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := strconv.Atoi(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	orders, err := h.client.ListOrdersByUser(c, &order.ListOrdersRequest{UserId: int32(userID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder returns a single order
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// Get user ID from context for authorization
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	order, err := h.client.GetOrder(c, &order.OrderID{Id: int32(id)})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Check if the order belongs to the current user
	userID, _ := strconv.Atoi(userIDStr.(string))
	if int(order.UserId) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized to view this order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var orderReq struct {
		Items []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, _ := strconv.Atoi(userIDStr.(string))

	// Create the order request
	newOrder := &order.Order{
		UserId: int32(userID),
		Status: "pending",
		Items:  make([]*order.OrderItem, 0, len(orderReq.Items)),
	}

	// Add items to the order
	for _, item := range orderReq.Items {
		newOrder.Items = append(newOrder.Items, &order.OrderItem{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	// Create the order
	createdOrder, err := h.client.CreateOrder(c, newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

// UpdateOrderStatus updates the status of an order
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrder, err := h.client.UpdateOrderStatus(c, &order.Order{
		Id:     int32(id),
		Status: statusUpdate.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedOrder)
}
