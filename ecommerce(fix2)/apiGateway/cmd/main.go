// package cmd
package main

import (
	"apiGateway/internal/delivery/handlers"
	grpcDelivery "apiGateway/internal/grpc"
	"apiGateway/internal/middleware"
	"log"
	_ "net/http"
	_ "net/http/httputil"
	_ "net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize gRPC clients
	userClient, err := grpcDelivery.NewUserClient("localhost:50053")
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	
	// Initialize Redis client
	redisClient := grpcDelivery.NewRedisClient()

	// Initialize Gin router
	r := gin.Default()

	// Add authentication middleware
	r.Use(middleware.Auth(userClient))

	// Register routes
	handlers.RegisterRoutes(r, userClient, redisClient)

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupServiceProxies настраивает проксирование запросов к микросервисам
func setupServiceProxies(r *gin.Engine) {
	// Create HTTP endpoints for inventory service
	inventoryHandler := createInventoryHandler()
	orderHandler := createOrderHandler()

	// REST API routes
	r.GET("/products", inventoryHandler.GetProducts)
	r.GET("/products/:id", inventoryHandler.GetProduct)
	r.POST("/products", inventoryHandler.CreateProduct)
	r.PUT("/products/:id", inventoryHandler.UpdateProduct)
	r.DELETE("/products/:id", inventoryHandler.DeleteProduct)

	r.GET("/orders", orderHandler.GetOrders)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.POST("/orders", orderHandler.CreateOrder)
	r.PUT("/orders/:id/status", orderHandler.UpdateOrderStatus)
}

// Implement these handler creators to connect to gRPC services
func createInventoryHandler() *handlers.InventoryHandler {
	// Connect to inventory gRPC service
	inventoryClient, err := grpcDelivery.NewInventoryClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to inventory service: %v", err)
	}
	return handlers.NewInventoryHandler(inventoryClient)
}

func createOrderHandler() *handlers.OrderHandler {
	// Connect to order gRPC service
	orderClient, err := grpcDelivery.NewOrderClient("localhost:50052")
	if err != nil {
		log.Fatalf("failed to connect to order service: %v", err)
	}
	return handlers.NewOrderHandler(orderClient)
}
