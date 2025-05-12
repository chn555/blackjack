/*
The blackjack executable is a gRPC server that serves
a blackjack service for managing decks of cards.

It listens on port 8081.

It expects a Deck Service to be available on port 8080.
*/
package main

import (
	"log/slog"
	"net"
	"os"

	blackjackImpl "github.com/chn555/blackjack/internal/blackjack"
	"github.com/chn555/blackjack/pkg/blackjack"
	blackjackPb "github.com/chn555/schemas/proto/blackjack/v1"
	deckPb "github.com/chn555/schemas/proto/deck/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	blackjackServer, err := blackjackImpl.NewServiceServer(store, deckClient)
	if err != nil {
		slog.Error("Error creating blackjack server", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
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
