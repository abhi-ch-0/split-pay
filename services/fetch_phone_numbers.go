package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AppService) FetchPhoneNumbers(ctx context.Context, _ *emptypb.Empty) (*pb.FetchPhoneNumbersResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `SELECT phone_number, is_searchable FROM backup_phone_numbers WHERE user_id = $1`
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []*pb.PhoneNumber
	for rows.Next() {
		var phone pb.PhoneNumber
		if err := rows.Scan(&phone.Contact, &phone.IsSearchable); err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, &phone)
	}

	return &pb.FetchPhoneNumbersResponse{PhoneNumbers: phoneNumbers}, nil
}
