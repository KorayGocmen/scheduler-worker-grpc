package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server holds the GRPC worker server instance.
type server struct{}

// StartJob
func (s *server) StartJob(ctx context.Context, r *pb.StartJobReq) (*pb.StartJobRes, error) {
	success, err, jobID := startScript(r.Command, r.Path)

	res := pb.StartJobRes{
		Success: success,
		Error:   err,
		JobID:   jobID,
	}
	return &res, nil
}

// StopJob
func (s *server) StopJob(ctx context.Context, r *pb.StopJobReq) (*pb.StopJobRes, error) {
	success, err := stopScript(r.JobID)

	res := pb.StopJobRes{
		Success: success,
		Error:   err,
	}
	return &res, nil
}

// QueryJob
func (s *server) QueryJob(ctx context.Context, r *pb.QueryJobReq) (*pb.QueryJobRes, error) {
	success, err, done := queryScript(r.JobID)

	res := pb.QueryJobRes{
		Success: success,
		Error:   err,
		Done:    done,
	}
	return &res, nil
}

// startGRPCServer starts the GRPC server for the worker.
// Scheduler can make grpc requests to this server to start,
// stop, query status of jobs etc.
func startGRPCServer() {
	lis, err := net.Listen("tcp", config.GRPCServer.Addr)
	if err != nil {
		fatal(fmt.Sprintf("failed to listen: %v", err))
	}

	var opts []grpc.ServerOption
	if config.GRPCServer.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(
			config.GRPCServer.CrtFile,
			config.GRPCServer.KeyFile,
		)
		if err != nil {
			fatal(fmt.Sprint("Failed to generate credentials", err))
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	log.Println("GRPC Server listening on", config.GRPCServer.Addr)

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
