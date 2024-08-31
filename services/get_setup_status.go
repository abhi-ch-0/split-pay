package services

import (
	"context"
	"fmt"
	pb "split-pay/generated"
	"split-pay/shared"
)

func (s *AppService) GetSetupStatus(ctx context.Context) (*pb.GetSetupStatusResponse, error) {
	userId, err := shared.GetUserIdFromToken(ctx)
	if err != nil {
		return nil, err
	}

	var hasDisplayName bool
	err = s.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users_display_name WHERE user_id = $1)", userId).Scan(&hasDisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to check display name: %v", err)
	}

	return &pb.GetSetupStatusResponse{
		IsSetupCompleted: hasDisplayName,
	}, nil
}
