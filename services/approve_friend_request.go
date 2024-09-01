package services

import (
	"context"
	"database/sql"
	pb "split-pay/generated"
	"split-pay/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *AppService) ApproveFriendRequest(ctx context.Context, req *pb.ApproveFriendRequestInput) (*pb.ApproveFriendRequestOutput, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateQuery := `UPDATE friend_requests SET status = 'approved' WHERE from_user_id = $1 AND to_user_id = $2 AND status = 'pending'`

	result, err := s.DB.Exec(updateQuery, req.FromUserId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve friend request: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "no pending request found or already approved")
	}

	insertQuery := `
		INSERT INTO friends (user1_id, user2_id)
		VALUES (LEAST($1, $2), GREATEST($1, $2))
		ON CONFLICT (user1_id, user2_id) DO NOTHING
	`
	_, err = tx.ExecContext(ctx, insertQuery, req.FromUserId, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add friendship: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &pb.ApproveFriendRequestOutput{
		Success: true,
	}, nil
}
