package main

import (
	"sync"

	"github.com/google/uuid"
)

var (
	workersMutex = &sync.Mutex{}
	workers      = make(map[string]*worker)
)

// worker holds the information about registered workers
// 		- id: uuid assigned when the worker first register.
// 		- addr: workers network address, later used to create
// 				grpc client to the worker
type worker struct {
	id   string
	addr string
}

// newWorker creates a new worker instance and adds
// the new worker to the map.
// Returns:
// 		- string: worker id
func newWorker(address string) string {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	workerID := uuid.New().String()
	workers[workerID] = &worker{
		id:   workerID,
		addr: address,
	}

	return workerID
}

// delWorker removes the worker with the given id
// from known workers map.
func delWorker(id string) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	if _, ok := workers[id]; ok {
		delete(workers, id)
	}
}
