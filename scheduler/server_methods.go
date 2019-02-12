package main

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/koraygocmen/gravitational/jobscheduler"
)

// RegisterWorker registers a new worker on the server.
func (s *server) RegisterWorker(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRes, error) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	workerID := uuid.New().String()
	workers[workerID] = &worker{
		id:   workerID,
		addr: r.Address,
	}

	res := pb.RegisterRes{
		Success: true,
		ID:      workerID,
	}

	return &res, nil
}
