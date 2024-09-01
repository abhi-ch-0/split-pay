package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) SendFriendRequest(ctx context.Context, req *pb.SendFriendRequestInput) (*pb.SendFriendRequestOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `INSERT INTO friend_requests (sender_id, receiver_id) VALUES ($1, $2)`
	_, err = s.DB.Exec(query, userId, req.RecipientId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send friend request: %v", err)
	}

	return &pb.SendFriendRequestOutput{
		Success: true,
	}, nil
}
