package services

import (
	"database/sql"
	pb "split-pay/generated"
)

type AppService struct {
	DB *sql.DB
	pb.UnimplementedSplitPayAppServiceServer
}
