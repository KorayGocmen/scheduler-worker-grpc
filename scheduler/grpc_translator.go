package main

import (
	"context"
	"errors"
	"time"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
)

// startJobOnWorker translates the http start request to grpc
// request on the workers.
// Returns:
// 		- string: job id
// 		- error: nil if no error
func startJobOnWorker(req apiStartJobReq) (string, error) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return "", errors.New("worker not found")
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return "", err
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
		return "", err
	}

	return r.JobID, nil
}

// stopJobOnWorker translates the http stop request to grpc
// request on the workers.
// Returns:
// 		- error: nil if no error
func stopJobOnWorker(req apiStopJobReq) error {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return errors.New("worker not found")
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewWorkerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stopJobReq := pb.StopJobReq{
		JobID: req.JobID,
	}

	if _, err := c.StopJob(ctx, &stopJobReq); err != nil {
		return err
	}

	return nil
}

// queryJobOnWorker translates the http query request to grpc
// request on the workers.
// Returns:
//		- bool: job status (true if job is done)
//		- bool: job error (true if job had an error)
// 		- string: job error text ("" if job error is false)
//		- error: nil if no error
func queryJobOnWorker(req apiQueryJobReq) (bool, bool, string, error) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	worker, ok := workers[req.WorkerID]
	if !ok {
		return false, false, "", errors.New("worker not found")
	}

	conn, err := grpc.Dial(worker.addr, grpc.WithInsecure())
	if err != nil {
		return false, false, "", err
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
		return false, false, "", err
	}

	return r.Done, r.Error, r.ErrorText, nil
}
