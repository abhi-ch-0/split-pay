package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) SetDisplayName(ctx context.Context, req *pb.SetDisplayNameRequest) (*pb.SetDisplayNameResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	if req.DisplayName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "display name cannot be empty")
	}

	query := `
        INSERT INTO users_display_name (user_id, display_name)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE
        SET display_name = EXCLUDED.display_name
    `
	_, err = s.DB.Exec(query, userId, req.DisplayName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set display name: %v", err)
	}

	return &pb.SetDisplayNameResponse{
		Success: true,
	}, nil
}
