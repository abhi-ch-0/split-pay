package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AppService) FetchGroups(ctx context.Context, _ *emptypb.Empty) (*pb.FetchGroupsOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `
		SELECT g.id, g.name
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.member_id = $1
	`
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch groups: %v", err)
	}
	defer rows.Close()

	var groups []*pb.Group
	for rows.Next() {
		var group pb.Group
		if err := rows.Scan(&group.GroupId, &group.Name); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan group: %v", err)
		}
		groups = append(groups, &group)
	}

	return &pb.FetchGroupsOutput{
		Groups: groups,
	}, nil
}
