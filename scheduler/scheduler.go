package main

import (
	"flag"
	"log"
	"net"
	"sync"

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
)

type server struct{}

var (
	workersMutex = &sync.Mutex{}
	workers      = make(map[string]*worker)
)

type worker struct {
	id   string
	addr string
}

func main() {

	go api()

	flag.Parse()
	lis, err := net.Listen("tcp", schedulerAddr)
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
	pb.RegisterSchedulerServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
