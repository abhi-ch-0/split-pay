package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) RemovePhoneNumber(ctx context.Context, req *pb.RemovePhoneNumberRequest) (*pb.RemovePhoneNumberResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `DELETE FROM backup_phone_numbers WHERE user_id = $1 AND phone_number = $2`
	result, err := s.DB.Exec(query, userId, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "failed to delete phone number: %v", err)
	}

	return &pb.RemovePhoneNumberResponse{Success: true}, nil
}
