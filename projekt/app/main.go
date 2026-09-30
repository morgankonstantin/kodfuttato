package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"code-executor/executorpb"
)

var version = "dev"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("nem sikerült a portot megnyitni: %v", err)
	}

	s := grpc.NewServer()
	executorpb.RegisterCodeExecutorServer(s, &server{})
	reflection.Register(s) // grpcurl proto fájl nélkül is tudja hívni

	log.Printf("code-executor %s fut a :%s porton", version, port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("a szerver leállt: %v", err)
	}
}
