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
	pb "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func RunGrpcServer(vocabManager *VocabManager, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration) {
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
	pb.RegisterVocabManagerServiceServer(s, &server{vocabs: vocabManager})

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

	// grpc does not have ctx-aware stop
	go func() {
		time.Sleep(shutdownTimeout)
		logger.Warn("timeout shutting down")
		os.Exit(1)
	}()

	s.GracefulStop()
	logger.Info("GRPC server stopped")
}

type server struct {
	pb.UnimplementedVocabManagerServiceServer
	vocabs *VocabManager
}

func (s *server) Ping(_ context.Context, _ *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Ok: true}, nil
}

func (s *server) GetAvailableVocabs(_ context.Context, _ *emptypb.Empty) (*pb.GetAvailableVocabsResponse, error) {
	return &pb.GetAvailableVocabsResponse{Names: s.vocabs.AvailableVocabs()}, nil
}

func (s *server) GetVocab(_ context.Context, req *pb.GetVocabRequest) (*pb.GetVocabResponse, error) {
	vocab := s.vocabs.Vocab(req.Name)
	return &pb.GetVocabResponse{PrimaryWords: vocab.PrimaryWords, RudeWords: vocab.RudeWords}, nil
}
