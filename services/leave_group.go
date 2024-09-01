package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) LeaveGroup(ctx context.Context, req *pb.LeaveGroupInput) (*pb.LeaveGroupOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `DELETE FROM group_members WHERE group_id = $1 AND member_id = $2`
	result, err := s.DB.Exec(query, req.GroupId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to leave group: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pb.LeaveGroupOutput{
			Success: false,
		}, status.Errorf(codes.InvalidArgument, "user is not a member of the group")
	}

	return &pb.LeaveGroupOutput{
		Success: true,
	}, nil
}
