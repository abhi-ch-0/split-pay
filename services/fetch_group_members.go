package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) FetchGroupMembers(ctx context.Context, req *pb.FetchGroupMembersInput) (*pb.FetchGroupMembersOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	var isMember bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND member_id = $2)`
	err = s.DB.QueryRow(checkQuery, req.GroupId, userId).Scan(&isMember)
	if err != nil || !isMember {
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of the group")
	}

	query := `
		SELECT u.id, u.username, u.display_name
		FROM users u
		JOIN group_members gm ON u.id = gm.member_id
		WHERE gm.group_id = $1
	`
	rows, err := s.DB.Query(query, req.GroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch group members: %v", err)
	}
	defer rows.Close()

	var members []*pb.GroupMember
	for rows.Next() {
		var member pb.GroupMember
		if err := rows.Scan(&member.UserId, &member.Username, &member.DisplayName); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to scan group member: %v", err)
		}
		members = append(members, &member)
	}

	return &pb.FetchGroupMembersOutput{
		Members: members,
	}, nil
}
