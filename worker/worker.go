package main

import (
	"context"
	"flag"
	"log"
	"net"
	"time"

	pb "github.com/koraygocmen/gravitational/jobscheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	schedulerAddr = "127.0.0.1:50000"
	workerAddr    = "127.0.0.1:30000"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")

	id string
)

func main() {
	// Set up a connection to the server and register
	conn, err := grpc.Dial(schedulerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSchedulerClient(conn)

	// Register here.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registerReq := pb.RegisterReq{
		Address: workerAddr,
	}
	r, err := c.RegisterWorker(ctx, &registerReq)
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Success: %t, ID: %s", r.Success, r.ID)

	// Starting the worker server
	flag.Parse()
	lis, err := net.Listen("tcp", workerAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = "server.pem"
		}
		if *keyFile == "" {
			*keyFile = "server.key"
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWorkerServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
