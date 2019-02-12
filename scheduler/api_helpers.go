package main

import (
	"context"
	"time"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
)

func startJobOnWorker(req rStartJobReq) (bool, string) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return false, "Worker not found."
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return false, err.Error()
	}
	defer conn.Close()
	c := pb.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	startJobReq := pb.StartJobReq{
		Command: req.Command,
		Path:    req.Path,
	}

	r, err := c.StartJob(ctx, &startJobReq)
	if err != nil {
		return false, err.Error()
	}

	return r.Success, r.Error
}
