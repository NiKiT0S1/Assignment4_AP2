package grpcDelivery

import (
	"apiGateway/internal/proto"
	"context"
	"google.golang.org/grpc"
	"log"
)

type UserClient struct {
	client proto.UserServiceClient
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // В реальной ситуации использовать grpc.WithTransportCredentials
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
		return nil, err
	}

	client := proto.NewUserServiceClient(conn)
	return &UserClient{client}, nil
}

// Authenticate перенаправляет запрос авторизации на gRPC сервис
func (u *UserClient) Authenticate(ctx context.Context, req *proto.AuthRequest) (*proto.UserResponse, error) {
	return u.client.Authenticate(ctx, req)
}

// Register перенаправляет запрос регистрации на gRPC сервис
func (u *UserClient) Register(username, password string) (*proto.UserResponse, error) {
	req := &proto.RegisterRequest{
		Username: username,
		Password: password,
	}
	return u.client.Register(context.Background(), req)
}

// GetProfile получает профиль пользователя по ID
func (u *UserClient) GetProfile(id int32) (*proto.UserResponse, error) {
	req := &proto.UserID{
		Id: id,
	}
	return u.client.GetProfile(context.Background(), req)
}

// UpdateProfile обновляет профиль пользователя
func (u *UserClient) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest) (*proto.UserResponse, error) {
	return u.client.UpdateProfile(ctx, req)
}
