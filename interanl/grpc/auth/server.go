package auth

import (
	"context"
	"errors"
	grpcAuthv1 "github.com/Xryak-Git/grpcAuthProto/gen/go/grpcAuth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpcAuth/interanl/services"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	Register(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userId int64) (isAdmin bool, err error)
}

type serverAPI struct {
	grpcAuthv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	grpcAuthv1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *grpcAuthv1.LoginRequest) (*grpcAuthv1.LoginResponse, error) {
	//TODO: add validation

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.AppId))

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &grpcAuthv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *grpcAuthv1.RegisterRequest) (*grpcAuthv1.RegisterResponse, error) {
	//TODO: add validation
	userID, err := s.auth.Register(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, "interanl error")
	}
	return &grpcAuthv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *grpcAuthv1.IsAdminRequest) (*grpcAuthv1.IsAdminResponse, error) {
	//TODO: add validation
	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &grpcAuthv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
