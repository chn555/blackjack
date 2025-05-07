package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/chn555/blackjack/pkg/ai"
	aiPb "github.com/chn555/schemas/proto/ai/v1"
)

type ServiceServer struct {
	aiPb.UnimplementedAiServiceServer
	aiServer *ai.AI
}

func NewServiceServer(aiServer *ai.AI) (*ServiceServer, error) {
	if aiServer == nil {
		return nil, fmt.Errorf("ai server is nil")
	}
	return &ServiceServer{aiServer: aiServer}, nil
}

func (s ServiceServer) PlayGame(ctx context.Context, request *aiPb.PlayGameRequest) (*aiPb.Empty, error) {
	err := s.aiServer.AddGame(request.GetGameId(), request.GetPlayerName(), s.protoToStrategy(request.GetStrategy()))
	if err != nil {
		return nil, fmt.Errorf("failed to add game: %w", err)
	}

	return &aiPb.Empty{}, nil
}

func (s ServiceServer) protoToStrategy(strategy string) ai.Strategy {
	switch {
	case strings.EqualFold(strategy, "dealer"):
		return &ai.DealerStrategy{}
	case strings.EqualFold(strategy, "greedy"):
		return &ai.GreedyStrategy{}
	default:
		return nil
	}
}
