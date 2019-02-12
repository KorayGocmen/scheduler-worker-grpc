package main

import (
	"context"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
)

// StartJob
func (s *server) StartJob(ctx context.Context, r *pb.StartJobReq) (*pb.StartJobRes, error) {
	success, err := startScript(r.Command, r.Path)

	res := pb.StartJobRes{
		Success: success,
		Error:   err,
	}
	return &res, nil
}

// StopJob
func (s *server) StopJob(ctx context.Context, r *pb.StopJobReq) (*pb.StopJobRes, error) {
	success, err := stopScript(r.Path)

	res := pb.StopJobRes{
		Success: success,
		Error:   err,
	}
	return &res, nil
}

// QueryJob
func (s *server) QueryJob(ctx context.Context, r *pb.QueryJobReq) (*pb.QueryJobRes, error) {
	success, err, done := queryScript(r.Path)

	res := pb.QueryJobRes{
		Success: success,
		Error:   err,
		Done:    done,
	}
	return &res, nil
}
