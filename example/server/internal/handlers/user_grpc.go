package handlers

import (
	"context"
	"fmt"

	"github.com/hoyci/gopher-chaos/example/pb"
	"github.com/hoyci/gopher-chaos/example/server/internal/repositories"
	"github.com/hoyci/gopher-chaos/example/server/internal/services"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	UseCase *services.UserUseCase
}

func (h *UserGRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := h.UseCase.CreateUser(req.Name, req.Email, req.Age)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: h.toProtoUser(user),
	}, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := h.UseCase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: h.toProtoUser(user),
	}, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := h.UseCase.UpdateByID(req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		User: h.toProtoUser(user),
	}, nil
}

func (h *UserGRPCHandler) ListUsers(req *pb.ListUserRequest, stream pb.UserService_ListUsersServer) error {
	for i := 0; i < int(req.Count); i++ {
		user := &repositories.User{
			ID:    fmt.Sprintf("%d", i),
			Name:  "Stream User",
			Email: "stream@chaos.com",
		}

		if err := stream.Send(&pb.UserResponse{User: h.toProtoUser(user)}); err != nil {
			return err
		}
	}
	return nil
}

func (h *UserGRPCHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.UseCase.DeleteByID(req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{
			Message: "error trying to delete user",
			Success: false,
		}, err
	}

	return &pb.DeleteUserResponse{
		Message: "user deleted successfully",
		Success: true,
	}, nil
}

func (h *UserGRPCHandler) toProtoUser(u *repositories.User) *pb.User {
	return &pb.User{
		Id:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Age:   u.Age,
	}
}
