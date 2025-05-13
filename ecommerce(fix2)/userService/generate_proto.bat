@echo off
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto
move proto\user.pb.go internal\delivery\grpc\pb\
move proto\user_grpc.pb.go internal\delivery\grpc\pb\ 