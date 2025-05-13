package grpc

import (
	"context"
	pb "inventoryService/internal/delivery/grpc/pb"
	"inventoryService/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	pb.UnimplementedInventoryServiceServer
	productUC domain.ProductUsecase
}

func NewInventoryHandler(productUC domain.ProductUsecase) *InventoryHandler {
	return &InventoryHandler{productUC: productUC}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	p := &domain.Product{
		Name: req.Name, Description: req.Description,
		Price: req.Price, Stock: int(req.Stock),
	}
	if err := h.productUC.Create(p); err != nil {
		return nil, status.Errorf(codes.Internal, "create failed: %v", err)
	}
	return toProto(p), nil
}

func (h *InventoryHandler) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.Product, error) {
	p, err := h.productUC.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "not found: %v", err)
	}
	return toProto(p), nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	p := &domain.Product{
		ID: int(req.Id), Name: req.Name,
		Description: req.Description, Price: req.Price, Stock: int(req.Stock),
	}
	if err := h.productUC.Update(p); err != nil {
		return nil, status.Errorf(codes.Internal, "update failed: %v", err)
	}
	return toProto(p), nil
}

func (h *InventoryHandler) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.Empty, error) {
	if err := h.productUC.Delete(int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}
	return &pb.Empty{}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, _ *pb.Empty) (*pb.ProductList, error) {
	products, err := h.productUC.List()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list failed: %v", err)
	}
	var res pb.ProductList
	for _, p := range products {
		res.Products = append(res.Products, toProto(&p))
	}
	return &res, nil
}

func toProto(p *domain.Product) *pb.Product {
	return &pb.Product{
		Id: int32(p.ID), Name: p.Name,
		Description: p.Description,
		Price:       p.Price, Stock: int32(p.Stock),
	}
}
