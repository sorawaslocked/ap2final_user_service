package grpc

import (
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/security"
	"github.com/sorawaslocked/ap2final_protos_gen/base"
	svc "github.com/sorawaslocked/ap2final_protos_gen/service/user"
	"github.com/sorawaslocked/ap2final_user_service/internal/adapter/grpc/dto"
	"github.com/sorawaslocked/ap2final_user_service/internal/model"
	"log/slog"
)

type UserServer struct {
	uc  UserUseCase
	log *slog.Logger
	svc.UnimplementedUserServiceServer
}

func NewUserServer(uc UserUseCase, log *slog.Logger) *UserServer {
	return &UserServer{
		uc:  uc,
		log: log,
	}
}

func (s *UserServer) Register(ctx context.Context, req *svc.RegisterRequest) (*svc.RegisterResponse, error) {
	const op = "grpc.UserServer.Register"

	log := s.log.With(slog.String("op", op))

	user := dto.ToUserFromRegisterRequest(req)

	registeredUser, err := s.uc.Register(ctx, user)
	if err != nil {
		logError(log, "register", err)

		return nil, dto.FromError(err)
	}

	return &svc.RegisterResponse{
		User: dto.FromUserToPb(registeredUser),
	}, nil
}

func (s *UserServer) Login(ctx context.Context, req *svc.LoginRequest) (*svc.LoginResponse, error) {
	const op = "grpc.UserServer.Login"

	log := s.log.With(slog.String("op", op))

	if req.Email == "" || req.Password == "" {
		err := dto.ErrMissingLoginCredentials
		logError(log, "login", err)

		return nil, dto.FromError(err)
	}

	user := dto.ToUserFromLoginRequest(req)

	token, err := s.uc.Login(ctx, user)
	if err != nil {
		logError(log, "login", err)

		return nil, dto.FromError(err)
	}

	return &svc.LoginResponse{
		Token: &base.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (s *UserServer) RefreshToken(ctx context.Context, req *svc.RefreshTokenRequest) (*svc.RefreshTokenResponse, error) {
	const op = "grpc.UserServer.RefreshToken"

	log := s.log.With(slog.String("op", op))

	token, err := s.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logError(log, "refresh token", err)

		return nil, dto.FromError(err)
	}

	return &svc.RefreshTokenResponse{
		Token: &base.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (s *UserServer) Get(ctx context.Context, req *svc.GetRequest) (*svc.GetResponse, error) {
	const op = "grpc.UserServer.Get"

	log := s.log.With(slog.String("op", op))

	token, ok := security.TokenFromCtx(ctx)
	if !ok {
		err := dto.ErrUnauthenticated
		logError(log, "get", err)

		return nil, dto.FromError(err)
	}

	user, err := s.uc.GetByID(
		ctx,
		model.Token{AccessToken: token},
		req.ID,
	)
	if err != nil {
		logError(log, "get", err)

		return nil, dto.FromError(err)
	}

	return &svc.GetResponse{
		User: dto.FromUserToPb(user),
	}, nil
}

func (s *UserServer) Update(ctx context.Context, req *svc.UpdateRequest) (*svc.UpdateResponse, error) {
	const op = "grpc.UserServer.Update"

	log := s.log.With(slog.String("op", op))

	id, update, credentialsUpdate, err := dto.ToUserUpdateFromUpdateRequest(req)
	if err != nil {
		logError(log, "update", err)

		return nil, dto.FromError(err)
	}

	token, ok := security.TokenFromCtx(ctx)
	if !ok {
		err := dto.ErrUnauthenticated
		logError(log, "update", err)

		return nil, dto.FromError(err)
	}

	updatedUser, err := s.uc.UpdateByID(ctx, model.Token{AccessToken: token}, id, credentialsUpdate, update)
	if err != nil {
		logError(log, "update", err)

		return nil, dto.FromError(err)
	}

	return &svc.UpdateResponse{
		User: dto.FromUserToPb(updatedUser),
	}, nil
}

func (s *UserServer) Delete(ctx context.Context, req *svc.DeleteRequest) (*svc.DeleteResponse, error) {
	const op = "grpc.UserServer.Delete"

	log := s.log.With(slog.String("op", op))

	token, ok := security.TokenFromCtx(ctx)
	if !ok {
		err := dto.ErrUnauthenticated
		logError(log, "delete", err)

		return nil, dto.FromError(err)
	}

	deletedUser, err := s.uc.DeleteByID(ctx, model.Token{AccessToken: token}, req.ID)
	if err != nil {
		logError(log, "delete", err)

		return nil, dto.FromError(err)
	}

	return &svc.DeleteResponse{
		User: dto.FromUserToPb(deletedUser),
	}, nil
}
