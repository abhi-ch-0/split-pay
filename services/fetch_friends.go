package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AppService) FetchFriends(ctx context.Context, _ *emptypb.Empty) (*pb.FetchFriendsOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `
    SELECT u.id, u.username, u.display_name
    FROM users u
	JOIN friendships f ON (f.user1_id = u.id OR f.user2_id = u.id)
	WHERE (f.user1_id = $1 OR f.user2_id = $1)
	`
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch friends: %v", err)
	}
	defer rows.Close()

	var friends []*pb.Friend
	for rows.Next() {
		var friend pb.Friend
		if err := rows.Scan(&friend.UserId, &friend.Username, &friend.DisplayName); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan friend: %v", err)
		}
		friends = append(friends, &friend)
	}

	return &pb.FetchFriendsOutput{Friends: friends}, nil

}
