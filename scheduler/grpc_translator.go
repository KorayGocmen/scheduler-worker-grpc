package main

import (
	"context"
	"time"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
)

func startJobOnWorker(req apiStartJobReq) (bool, string, string) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return false, "Worker not found.", ""
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return false, err.Error(), ""
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
		return false, err.Error(), ""
	}

	return r.Success, r.Error, r.JobID
}

func stopJobOnWorker(req apiStopJobReq) (bool, string) {
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

	stopJobReq := pb.StopJobReq{
		JobID: req.JobID,
	}

	r, err := c.StopJob(ctx, &stopJobReq)
	if err != nil {
		return false, err.Error()
	}

	return r.Success, r.Error
}

func queryJobOnWorker(req apiQueryJobReq) (bool, string, bool) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return false, "Worker not found.", false
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return false, err.Error(), false
	}
	defer conn.Close()
	c := pb.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	queryJobReq := pb.QueryJobReq{
		JobID: req.JobID,
	}

	r, err := c.QueryJob(ctx, &queryJobReq)
	if err != nil {
		return false, err.Error(), false
	}

	return r.Success, r.Error, r.Done
}
