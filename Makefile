# Makefile

SCHEDULER_PATH = scheduler
SCHEDULER_BINARY = scheduler

WORKER_PATH = worker
WORKER_BINARY = worker

PROTO_PATH = jobscheduler
PROTO_NAME = job_scheduler.proto
PROTO_PLUGIN = jobscheduler

.PHONY: all build scheduler worker clean

all: clean build

build: build_scheduler build_worker

scheduler: build_scheduler run_scheduler

worker: build_worker run_worker

build_proto:
	protoc -I $(PROTO_PATH)/ $(PROTO_PATH)/$(PROTO_NAME) --go_out=plugins=grpc:$(PROTO_PLUGIN)

# Scheduler
build_scheduler:
	go build -o $(SCHEDULER_PATH)/$(SCHEDULER_BINARY) $(SCHEDULER_PATH)/*.go

run_scheduler:
	$(SCHEDULER_PATH)/$(SCHEDULER_BINARY)

# Worker
build_worker:
	go build -o $(WORKER_PATH)/$(WORKER_BINARY) $(WORKER_PATH)/*.go

run_worker:
	$(WORKER_PATH)/$(WORKER_BINARY)

clean: 
	rm -f *.out
	rm -f $(SCHEDULER_PATH)/$(SCHEDULER_BINARY)
	rm -f $(WORKER_PATH)/$(WORKER_BINARY)
