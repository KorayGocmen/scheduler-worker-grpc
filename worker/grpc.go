package main

import (
	"context"
	"log"
	"time"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
)

func registerWorker() {
	conn, err := grpc.Dial(config.Scheduler.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registerReq := pb.RegisterReq{
		Address: config.GRPCServer.Addr,
	}
	r, err := c.RegisterWorker(ctx, &registerReq)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}

	workerID = r.ID
	log.Printf("Registered OK: %t, ID: %s", r.Success, r.ID)
}

func deregisterWorker() {
	conn, err := grpc.Dial(config.Scheduler.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	deregisterReq := pb.DeregisterReq{
		ID: workerID,
	}
	r, err := c.DeregisterWorker(ctx, &deregisterReq)
	if err != nil {
		log.Fatalf("could not deregister: %v", err)
	}

	log.Printf("Deregistered OK: %t", r.Success)
}
