package main

import (
	"log"
	"net"
	grpcDelivery "userService/internal/delivery/grpc"
	pb "userService/internal/delivery/grpc/pb"
	"userService/internal/repository"
	"userService/internal/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=0000 dbname=ecommerce sslmode=disable")
	if err != nil {
		log.Fatalln("failed to connect to DB:", err)
	}

	userRepo := repository.NewUserRepo(db)
	userUC := usecase.NewUserUsecase(userRepo)
	handler := grpcDelivery.NewUserHandler(userUC)

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, handler)

	log.Println("UserService gRPC started on port 50053")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
