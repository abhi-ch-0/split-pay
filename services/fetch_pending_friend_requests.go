package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AppService) FetchPendingFriendRequests(ctx context.Context, _ *emptypb.Empty) (*pb.FetchPendingFriendRequestsOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `
    SELECT fr.from_user_id, u.username, u.display_name
    FROM friend_requests fr
    JOIN users u ON fr.from_user_id = u.id
    WHERE fr.to_user_id = $1 AND fr.status = 'pending'
    `
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch pending friend requests: %v", err)
	}
	defer rows.Close()

	var pendingFriendRequests []*pb.PendingFriendRequest
	for rows.Next() {
		var req pb.PendingFriendRequest
		if err := rows.Scan(&req.FromUserId, &req.FromUsername, &req.DisplayName); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan pending friend request: %v", err)
		}
		pendingFriendRequests = append(pendingFriendRequests, &req)
	}

	return &pb.FetchPendingFriendRequestsOutput{PendingFriendRequests: pendingFriendRequests}, nil
}
