package services

import (
	"context"
	"database/sql"
	"os"
	pb "split-pay/generated"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *AppService) LogIn(ctx context.Context, req *pb.LogInInput) (*pb.LogInOutput, error) {
	var userId string
	var storedHashedPassword string
	err := s.DB.QueryRowContext(ctx, "SELECT id, password_hash FROM users WHERE username = $1", req.Username).Scan(&userId, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LogInOutput{
				StatusCode: 401,
				Message:    "Invalid username",
			}, nil
		}
		return &pb.LogInOutput{
			StatusCode: 500,
			Message:    "Internal server error",
		}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.Password))
	if err != nil {
		return &pb.LogInOutput{
			StatusCode: 401,
			Message:    "Invalid password",
		}, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return &pb.LogInOutput{
			StatusCode: 500,
			Message:    "Error while creating token",
		}, err
	}

	return &pb.LogInOutput{
		StatusCode: 200,
		UserId:     userId,
		Token:      tokenString,
		Message:    "Login successful",
	}, nil
}
