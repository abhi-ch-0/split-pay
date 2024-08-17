package services

import (
	"context"
	pb "split-pay/generated"

	"golang.org/x/crypto/bcrypt"
)

func (s *AppService) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var doesUsernameExist bool
	err := s.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&doesUsernameExist)
	if err != nil {
		return &pb.SignUpResponse{
			StatusCode: 500,
			Message:    "Internal server error",
		}, err
	}
	if doesUsernameExist {
		return &pb.SignUpResponse{
			StatusCode: 400,
			Message:    "Username already exists",
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.SignUpResponse{
			StatusCode: 500,
			Message:    "Error while hashing password",
		}, err
	}

	var userId string
	err = s.DB.QueryRowContext(ctx, `
		INSERT INTO users (username, password_hash) 
		VALUES ($1, $2) 
		RETURNING id`,
		req.Username, string(hashedPassword),
	).Scan(&userId)

	if err != nil {
		return &pb.SignUpResponse{
			StatusCode: 500,
			Message:    "Error while creating user",
		}, err
	}

	token := "generated-jwt-token"

	return &pb.SignUpResponse{
		StatusCode: 200,
		UserId:     userId,
		Token:      token,
		Message:    "User created successfully",
	}, nil
}
