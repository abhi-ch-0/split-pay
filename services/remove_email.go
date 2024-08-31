package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) RemoveEmail(ctx context.Context, req *pb.RemoveEmailRequest) (*pb.RemoveEmailResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `DELETE FROM backup_emails WHERE user_id = $1 AND email = $2`
	result, err := s.DB.Exec(query, userId, req.Email)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "Email not found")
	}

	return &pb.RemoveEmailResponse{Success: true}, nil
}
