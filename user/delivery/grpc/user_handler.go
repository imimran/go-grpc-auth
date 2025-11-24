package grpc

import (
	"context"
	"github.com/imimran/go-grpc-auth/user/usecase"
	pb "github.com/imimran/go-grpc-auth/proto"
	transformer "github.com/imimran/go-grpc-auth/user/transformer/grpc"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user, err := h.userUsecase.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}
	return transformer.ToProtoUser(user), nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.UserId) (*pb.User, error) {
	user, err := h.userUsecase.Get(req.Id)
	if err != nil {
		return nil, err
	}
	return transformer.ToProtoUser(user), nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := h.userUsecase.Update(req.Id, req.Email, req.Password, req.FullName)
	if err != nil {
		return nil, err
	}
	return transformer.ToProtoUser(user), nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.UserId) (*pb.Empty, error) {
	err := h.userUsecase.Delete(req.Id)
	return &pb.Empty{}, err
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.Empty) (*pb.UserListResponse, error) {
    users, err := h.userUsecase.List()
    if err != nil {
        return nil, err
    }

    resp := &pb.UserListResponse{
        Users: transformer.ToProtoUserList(users),
    }

    return resp, nil
}


func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, nil
}
