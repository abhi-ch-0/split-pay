package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) CreateGroup(ctx context.Context, req *pb.CreateGroupInput) (*pb.CreateGroupOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `INSERT INTO groups (name, created_by) VALUES ($1, $2) RETURNING id`
	var groupId int32
	err = s.DB.QueryRow(query, req.GroupName, userId).Scan(&groupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create group: %v", err)
	}

	return &pb.CreateGroupOutput{
		Success: true,
		GroupId: groupId,
	}, nil
}
