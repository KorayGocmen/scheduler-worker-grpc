package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/koraygocmen/scheduler-worker-grpc/jobscheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server holds the GRPC worker server instance.
type server struct{}

// StartJob starts a new job with the given command and the path
// Command can be any exectuable command on the worker and the path
// is the relative path of the script.
func (s *server) StartJob(ctx context.Context, r *pb.StartJobReq) (*pb.StartJobRes, error) {
	jobID, err := startScript(r.Command, r.Path)
	if err != nil {
		return nil, err
	}

	res := pb.StartJobRes{
		JobID: jobID,
	}

	return &res, nil
}

// StopJob stops a running job with the given job id.
func (s *server) StopJob(ctx context.Context, r *pb.StopJobReq) (*pb.StopJobRes, error) {
	if err := stopScript(r.JobID); err != nil {
		return nil, err
	}

	return &pb.StopJobRes{}, nil
}

// QueryJob returns the status of job with the given job id.
// The status of the job is inside the `Done` variable in response
// and it specifies if the job is still running (true), or stopped (false).
func (s *server) QueryJob(ctx context.Context, r *pb.QueryJobReq) (*pb.QueryJobRes, error) {
	jobDone, jobError, jobErrorText, err := queryScript(r.JobID)
	if err != nil {
		return nil, err
	}

	res := pb.QueryJobRes{
		Done:      jobDone,
		Error:     jobError,
		ErrorText: jobErrorText,
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
