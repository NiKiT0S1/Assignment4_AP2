package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	// You'll need to create these proto imports
	"apiGateway/internal/proto/inventory"
)

// InventoryClient interface to make testing easier
type InventoryClient interface {
	CreateProduct(ctx context.Context, product *inventory.Product, opts ...grpc.CallOption) (*inventory.Product, error)
	GetProduct(ctx context.Context, id *inventory.ProductID, opts ...grpc.CallOption) (*inventory.Product, error)
	UpdateProduct(ctx context.Context, product *inventory.Product, opts ...grpc.CallOption) (*inventory.Product, error)
	DeleteProduct(ctx context.Context, id *inventory.ProductID, opts ...grpc.CallOption) (*inventory.Empty, error)
	ListProducts(ctx context.Context, empty *inventory.Empty, opts ...grpc.CallOption) (*inventory.ProductList, error)
}

// InventoryHandler handles HTTP requests for inventory service
type InventoryHandler struct {
	client InventoryClient
}

// NewInventoryHandler creates a new inventory handler
func NewInventoryHandler(client InventoryClient) *InventoryHandler {
	return &InventoryHandler{client: client}
}

// GetProducts returns all products
func (h *InventoryHandler) GetProducts(c *gin.Context) {
	products, err := h.client.ListProducts(c, &inventory.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct returns a single product
func (h *InventoryHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.client.GetProduct(c, &inventory.ProductID{Id: int32(id)})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var product inventory.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := h.client.CreateProduct(c, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct updates an existing product
func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var product inventory.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.Id = int32(id)

	updatedProduct, err := h.client.UpdateProduct(c, &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes a product
func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	_, err = h.client.DeleteProduct(c, &inventory.ProductID{Id: int32(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
