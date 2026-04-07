package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/xolra0d/alias-online/shared/proto/room_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RunGrpcServer(roomManager *RoomManager, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration) {
	const op = "main.RunGrpcServer"

	lis, err := net.Listen("tcp", runningAddr)
	if err != nil {
		logger.Error("failed to listen", "addr", runningAddr, "err", err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	pb.RegisterRoomManagerServiceServer(s, &server{
		roomManager: roomManager,
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
	pb.UnimplementedRoomManagerServiceServer
	roomManager *RoomManager
	logger      *slog.Logger

	SetWorkerActiveTimeout time.Duration
}

func (s *server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Ok: true}, nil
}

func (s *server) PingWorker(ctx context.Context, req *pb.PingWorkerRequest) (*emptypb.Empty, error) {
	var lastError error = nil

	err := s.roomManager.SetWorkerActive(ctx, req.GetWorker())

	if err != nil {
		s.logger.Error("failed to prolong worker", "worker", req.GetWorker(), "err", err)
		lastError = err
	}
	for _, roomId := range req.LoadedRooms {
		err = s.roomManager.ProlongRoom(ctx, roomId, req.GetWorker())
		if err != nil {
			s.logger.Error("failed to prolong room", "room", roomId, "err", err)
			lastError = err
		}
	}

	if lastError != nil {
		return nil, status.Error(codes.Internal, lastError.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *server) GetRoomWorker(ctx context.Context, req *pb.GetRoomWorkerRequest) (*pb.GetRoomWorkerResponse, error) {
	worker, err := s.roomManager.FindBestWorker(ctx, req.RoomId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetRoomWorkerResponse{Worker: worker}, nil
}

func (s *server) RegisterRoom(ctx context.Context, req *pb.RegisterRoomRequest) (*pb.RegisterRoomResponse, error) {
	worker, err := s.roomManager.RegisterRoom(ctx, req.RoomId, req.GetWorker())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RegisterRoomResponse{Worker: worker}, nil
}

func (s *server) ProlongRoom(ctx context.Context, req *pb.ProlongRoomRequest) (*emptypb.Empty, error) {
	err := s.roomManager.ProlongRoom(ctx, req.RoomId, req.Worker)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *server) ReleaseRoom(ctx context.Context, req *pb.ReleaseRoomRequest) (*emptypb.Empty, error) {
	err := s.roomManager.ReleaseRoom(ctx, req.RoomId, req.Worker)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
