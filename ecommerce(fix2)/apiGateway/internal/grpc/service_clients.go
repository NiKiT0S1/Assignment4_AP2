package grpcDelivery

import (
	"apiGateway/internal/proto/inventory"
	"apiGateway/internal/proto/order"
	"log"

	"google.golang.org/grpc"
)

// NewInventoryClient creates a new gRPC client for inventory service
func NewInventoryClient(address string) (inventory.InventoryServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // In production, use grpc.WithTransportCredentials
	if err != nil {
		log.Printf("failed to connect to inventory service: %v", err)
		return nil, err
	}

	client := inventory.NewInventoryServiceClient(conn)
	return client, nil
}

// NewOrderClient creates a new gRPC client for order service
func NewOrderClient(address string) (order.OrderServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // In production, use grpc.WithTransportCredentials
	if err != nil {
		log.Printf("failed to connect to order service: %v", err)
		return nil, err
	}

	client := order.NewOrderServiceClient(conn)
	return client, nil
}
