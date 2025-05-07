package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	aiImpl "github.com/chn555/blackjack/internal/ai"
	"github.com/chn555/blackjack/pkg/ai"
	aiPb "github.com/chn555/schemas/proto/ai/v1"
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Create a connection to the blackjack service
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("Error connecting to blackjack service", err)
		os.Exit(1)
	}
	defer conn.Close()

	blackjackClient := blackjackPb.NewBlackjackServiceClient(conn)

	// Create the blackjack server with the deck client
	a := ai.NewAI(blackjackClient)
	a.Start()

	aiServer, err := aiImpl.NewServiceServer(a)
	if err != nil {
		slog.Error("Error creating ai server", err)
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
	aiPb.RegisterAiServiceServer(grpcServer, aiServer)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		slog.Error("Error creating listener", err)
		os.Exit(1)
	}

	slog.Info("Listening on port 8082")
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
