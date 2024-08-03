package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "split-pay/generated"
	"split-pay/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	InitDB()
	defer CloseDB()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &services.GreeterService{})

	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		log.Println("gRPC server listening on port 50051")
	}()

	<-stop
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
