package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) DeclineFriendRequest(ctx context.Context, req *pb.DeclineFriendRequestInput) (*pb.DeclineFriendRequestOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `UPDATE friend_requests SET status = 'declined' WHERE from_user_id = $1 AND to_user_id = $2 AND status = 'pending'`

	result, err := s.DB.Exec(query, req.FromUserId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decline friend request: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "no pending request found or already declined")
	}

	return &pb.DeclineFriendRequestOutput{
		Success: true,
	}, nil
}
