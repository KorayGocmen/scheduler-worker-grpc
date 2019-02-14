package main

import (
	"sync"
)

var (
	workersMutex = &sync.Mutex{}
	workers      = make(map[string]*worker)
)

type worker struct {
	id   string
	addr string
}

func newWorker(address string) string {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	// TODO: Change this back
	// workerID := uuid.New().String()
	workerID := "test_worker"
	workers[workerID] = &worker{
		id:   workerID,
		addr: address,
	}

	return workerID
}

func delWorker(id string) {
	workersMutex.Lock()
	defer workersMutex.Unlock()

	if _, ok := workers[id]; ok {
		delete(workers, id)
	}
}
