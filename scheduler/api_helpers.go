package main

// apiStartJobReq expected API payload for `/start`
type apiStartJobReq struct {
	Command  string `json:"command"`
	Path     string `json:"path"`
	WorkerID string `json:"worker_id"`
}

// apiStartJobRes returned API payload for `/start`
type apiStartJobRes struct {
	JobID string `json:"job_id"`
}

// apiStopJobReq expected API payload for `/stop`
type apiStopJobReq struct {
	JobID    string `json:"job_id"`
	WorkerID string `json:"worker_id"`
}

// apiStopJobRes returned API payload for `/stop`
type apiStopJobRes struct {
	Success bool `json:"success"`
}

// apiQueryJobReq expected API payload for `/query`
type apiQueryJobReq struct {
	JobID    string `json:"job_id"`
	WorkerID string `json:"worker_id"`
}

// apiQueryJobRes returned API payload for `/query`
type apiQueryJobRes struct {
	Done bool `json:"done"`
}
