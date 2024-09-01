package services

import (
	"context"
	"os"
	pb "split-pay/generated"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *AppService) SignUp(ctx context.Context, req *pb.SignUpInput) (*pb.SignUpOutput, error) {
	var doesUsernameExist bool
	err := s.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&doesUsernameExist)
	if err != nil {
		return &pb.SignUpOutput{
			StatusCode: 500,
			Message:    "Internal server error",
		}, err
	}
	if doesUsernameExist {
		return &pb.SignUpOutput{
			StatusCode: 400,
			Message:    "Username already exists",
		}, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.SignUpOutput{
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
		return &pb.SignUpOutput{
			StatusCode: 500,
			Message:    "Error while creating user",
		}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return &pb.SignUpOutput{
			StatusCode: 500,
			Message:    "Error while creating token",
		}, err
	}

	return &pb.SignUpOutput{
		StatusCode: 200,
		UserId:     userId,
		Token:      tokenString,
		Message:    "User created successfully",
	}, nil
}
