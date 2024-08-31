package services

import (
	"context"
	"regexp"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) AddPhoneNumber(ctx context.Context, req *pb.AddPhoneNumberRequest) (*pb.AddPhoneNumberResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	phonePattern := `^\+?[\d\s\-()]{10,15}$`
	re := regexp.MustCompile(phonePattern)
	isPhoneNumberValid := re.MatchString(req.PhoneNumber.Contact)
	if !isPhoneNumberValid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid phone number")
	}

	query := `
        INSERT INTO backup_phone_numbers (user_id, phone_number, is_searchable)
        VALUES ($1, $2, $3)
    `
	_, err = s.DB.Exec(query, userId, req.PhoneNumber.Contact, req.PhoneNumber.IsSearchable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set phone number: %v", err)
	}

	return &pb.AddPhoneNumberResponse{
		Success: true,
	}, nil
}
