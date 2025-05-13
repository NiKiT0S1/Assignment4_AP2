package grpc

import (
	"context"
	pb "userService/internal/delivery/grpc/pb"
	"userService/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	uc domain.UserUsecase
}

func NewUserHandler(uc domain.UserUsecase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	u, err := h.uc.Register(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}
	return &pb.UserResponse{Id: int32(u.ID), Username: u.Username}, nil
}

func (h *UserHandler) Authenticate(ctx context.Context, req *pb.AuthRequest) (*pb.UserResponse, error) {
	u, err := h.uc.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	return &pb.UserResponse{Id: int32(u.ID), Username: u.Username}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	u, err := h.uc.GetProfile(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}
	return &pb.UserResponse{Id: int32(u.ID), Username: u.Username}, nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UserResponse, error) {
	u, err := h.uc.UpdateProfile(int(req.Id), req.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	return &pb.UserResponse{Id: int32(u.ID), Username: u.Username}, nil
}
