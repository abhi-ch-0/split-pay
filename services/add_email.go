package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) AddEmail(ctx context.Context, req *pb.AddEmailInput) (*pb.AddEmailOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `INSERT INTO backup_emails (user_id, email, is_searchable) VALUES ($1, $2, $3)`
	_, err = s.DB.Exec(query, userId, req.Email.Address, req.Email.IsSearchable)
	if err != nil {
		return nil, err
	}

	return &pb.AddEmailOutput{Success: true}, nil
}
