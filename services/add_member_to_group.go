package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) AddMemberToGroup(ctx context.Context, req *pb.AddMemberToGroupInput) (*pb.AddMemberToGroupOutput, error) {
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

	query := `INSERT INTO group_members (group_id, member_id, added_by) VALUES ($1, $2, $3)`
	_, err = s.DB.Exec(query, req.GroupId, req.NewMemberId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add member to group: %v", err)
	}

	return &pb.AddMemberToGroupOutput{
		Success: true,
	}, nil
}
