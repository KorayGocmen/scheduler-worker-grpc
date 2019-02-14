package main

type rStartJobReq struct {
	Command  string `json:"command"`
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

type rStopJobReq struct {
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

type rQueryJobReq struct {
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

type rStreamJobReq struct {
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}
