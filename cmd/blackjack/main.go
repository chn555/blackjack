package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	blackjackImpl "github.com/chn555/blackjack/internal/blackjack"
	"github.com/chn555/blackjack/pkg/blackjack"
	aiPb "github.com/chn555/schemas/proto/ai/v1"
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func main() {
	store := blackjack.NewInMemoryGameStore()

	// Create a connection to the deck service
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error connecting to deck service", err)
		os.Exit(1)
	}
	defer conn.Close()

	deckClient := deckPb.NewDeckServiceClient(conn)
	// Create a connection to the deck service
	aiConn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error connecting to deck service", err)
		os.Exit(1)
	}
	defer aiConn.Close()
	aiClient := aiPb.NewAiServiceClient(aiConn)
	// Create the blackjack server with the deck client
	blackjackServer, err := blackjackImpl.NewServiceServer(store, deckClient, aiClient)
	if err != nil {
		slog.Error("Error creating blackjack server", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		),
	)
	blackjackPb.RegisterBlackjackServiceServer(grpcServer, blackjackServer)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		slog.Error("Error creating listener", err)
		os.Exit(1)
	}

	slog.Info("Listening on port 8081")
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("Error serving grpc server", err)
		os.Exit(1)
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
