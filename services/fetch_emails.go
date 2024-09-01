package services

import (
	"context"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *AppService) FetchEmails(ctx context.Context, _ *emptypb.Empty) (*pb.FetchEmailsOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	query := `SELECT email, is_searchable FROM backup_emails WHERE user_id = $1`
	rows, err := s.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []*pb.Email
	for rows.Next() {
		var email pb.Email
		if err := rows.Scan(&email.Address, &email.IsSearchable); err != nil {
			return nil, err
		}
		emails = append(emails, &email)
	}

	return &pb.FetchEmailsOutput{Emails: emails}, nil
}
