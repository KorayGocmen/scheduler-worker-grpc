package main

import (
	"context"
	"log"
	"net"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server holds the GRPC scheduler server instance.
type server struct{}

// RegisterWorker registers a new worker on the server.
// Workers call this method when they are coming online.
func (s *server) RegisterWorker(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRes, error) {

	workerID := newWorker(r.Address)

	res := pb.RegisterRes{
		Success:  true,
		WorkerID: workerID,
	}

	log.Printf("New worker with ID: %s is online\n", workerID)
	return &res, nil
}

// DeregisterWorker deregisters an existing worker on the server.
// Workers call this method when they are going offline.
func (s *server) DeregisterWorker(ctx context.Context, r *pb.DeregisterReq) (*pb.DeregisterRes, error) {

	delWorker(r.WorkerID)

	res := pb.DeregisterRes{
		Success: true,
	}

	log.Printf("Worker with ID: %s is offline\n", r.WorkerID)
	return &res, nil
}

// startGRPCServer starts a scheduler server instance on the address specified
// by the config.GRPCServer.Addr, if the config.GRPCServer.UseTLS is true, the
// GRPC server will start with TLS with the key and crt file speficied in config.
func startGRPCServer() {
	lis, err := net.Listen("tcp", config.GRPCServer.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start with TLS if option is set in the config.
	var opts []grpc.ServerOption
	if config.GRPCServer.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(
			config.GRPCServer.CrtFile,
			config.GRPCServer.KeyFile,
		)
		if err != nil {
			log.Fatalln("Failed to generate credentials", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	log.Println("GRPC Server listening on", config.GRPCServer.Addr)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSchedulerServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
