package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/middleware"
	pb "github.com/xolra0d/alias-online/shared/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RunGrpcServer(secrets *Secrets, postgres *Postgres, logger *slog.Logger, addAccountTimeout, findAccountTimeout, jwtCookieTimeout time.Duration, runningAddr string, shutdownTimeout time.Duration) {
	lis, err := net.Listen("tcp", runningAddr)
	if err != nil {
		logger.Error("failed to listen", "addr", runningAddr, "err", err)
		os.Exit(1)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingUnaryInterceptor(logger),
		),
	)
	pb.RegisterAuthServiceServer(s, &server{
		secrets:  secrets,
		postgres: postgres,
		logger:   logger,

		AddAccountTimeout:  addAccountTimeout,
		FindAccountTimeout: findAccountTimeout,
		JWTCookieTimeout:   jwtCookieTimeout,
	})

	shutdown := make(chan os.Signal, 1)

	go func() {
		logger.Info("starting GRPC server", "addr", runningAddr)
		if err := s.Serve(lis); err != nil {
			logger.Error("failed to serve: %v", err)
		}
	}()

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	logger.Info("shutdown initiated")

	go func() {
		time.Sleep(shutdownTimeout)
		logger.Warn("timeout shutting down")
		os.Exit(1)
	}()

	s.GracefulStop()
	logger.Info("GRPC server stopped")
}

type server struct {
	pb.UnimplementedAuthServiceServer
	secrets  *Secrets
	postgres *Postgres
	logger   *slog.Logger

	AddAccountTimeout  time.Duration
	FindAccountTimeout time.Duration
	JWTCookieTimeout   time.Duration
}

func (s *server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Ok: true}, nil
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := ValidateForRegister(req.Name, req.Login, req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	hashed := s.secrets.hashSecret(req.Password)
	ctx, cancel := context.WithTimeout(ctx, s.AddAccountTimeout)
	err := s.postgres.AddAccount(ctx, req.Name, req.Login, hashed)
	cancel()
	if err != nil {
		switch err.SvcError {
		case ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		case ErrInternal:
			return nil, status.Error(codes.Internal, "internal server error")
		default:
			return nil, status.Error(codes.Unknown, "unknown error")
		}
	}

	exp := time.Now().Add(s.JWTCookieTimeout)
	token, err2 := s.secrets.NewJWT(req.Login, exp)
	if err2 != nil {
		return nil, status.Error(codes.Internal, err2.Error())
	}

	return &pb.RegisterResponse{Token: token, Exp: exp.Unix()}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := ValidateForLogin(req.Login, req.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, s.AddAccountTimeout)
	hashed, err := s.postgres.FindAccount(ctx, req.Login)
	cancel()
	if err != nil {
		switch err.SvcError {
		case ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "user not found")
		case ErrInternal:
			return nil, status.Error(codes.Internal, "internal server error")
		default:
			return nil, status.Error(codes.Unknown, "unknown error")
		}
	}
	err = s.secrets.VerifyPassword(req.Password, hashed)
	if err != nil {
		switch err.SvcError {
		case ErrWrongCredentials:
			return nil, status.Error(codes.Unauthenticated, "wrong credentials")
		case ErrInternal:
			return nil, status.Error(codes.Internal, "internal server error")
		default:
			return nil, status.Error(codes.Unknown, "unknown error")
		}
	}

	exp := time.Now().Add(s.JWTCookieTimeout)
	token, err2 := s.secrets.NewJWT(req.Login, exp)
	if err2 != nil {
		return nil, status.Error(codes.Internal, err2.Error())
	}

	return &pb.LoginResponse{Token: token, Exp: exp.Unix()}, nil
}
