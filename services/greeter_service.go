package services

import (
	"context"
	pb "split-pay/generated"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	reply := &pb.HelloReply{
		Message: "Hello, " + req.Name,
	}
	return reply, nil
}
