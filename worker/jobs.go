package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/google/uuid"
)

// jobsMutex is the lock to access jobs map.
// jobs is the map that holds current/past jobs.
// 		- key: job id
// 		- value: pointer to the created job object.
var (
	jobsMutex = &sync.Mutex{}
	jobs      = make(map[string]*job)
)

// job holds information about the ongoing or past jobs,
// that were triggered by the scheduler.
// 		- id: UUID assigned by the worker and sent back to the scheduler.
// 		- command: command which the scheduler run the job with
// 		- path: path to the job file/executable sent by the scheduler.
// 		- outFilePath: file path to where the output of the job will be piped.
// 		- cmd: pointer to the cmd.Exec command to get job status etc.
// 		- done: whether if job is done (default false)
//    - err: error while running the job (default nil)
type job struct {
	id          string
	command     string
	path        string
	outFilePath string
	cmd         *exec.Cmd
	done        bool
	err         error
}

// startScript start a new job.
// Returns:
//		- string: job id
//		- error: nil if no error
func startScript(command, path string) (string, error) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	jobID := uuid.New().String()
	outFilePath := fmt.Sprintf("%s.out", jobID)

	outfile, err := os.Create(outFilePath)
	if err != nil {
		return "", err
	}
	defer outfile.Close()

	cmd := exec.Command(command, path)
	cmd.Stdout = outfile

	if err = cmd.Start(); err != nil {
		return "", err
	}

	newJob := job{
		id:          jobID,
		command:     command,
		path:        path,
		outFilePath: outFilePath,
		cmd:         cmd,
		done:        false,
		err:         nil,
	}
	jobs[jobID] = &newJob

	// Get the status of the job async.
	go func() {
		if err := cmd.Wait(); err != nil {
			newJob.err = err
		}
		newJob.done = true
	}()

	return jobID, nil
}

// stopScript stop a running job.
// Returns:
//		- error: nil if no error
func stopScript(jobID string) error {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, found := jobs[jobID]
	if !found {
		return errors.New("job not found")
	}

	if job.done {
		return nil
	}

	if err := job.cmd.Process.Kill(); err != nil {
		return err
	}

	return nil
}

// queryScript check if job is done or not.
// Returns:
//		- bool: job status (true if job is done)
//		- bool: job error (true if job had an error)
// 		- string: job error text ("" if job error is false)
//		- error: nil if no error
func queryScript(jobID string) (bool, bool, string, error) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, found := jobs[jobID]
	if !found {
		return false, false, "", errors.New("job not found")
	}

	var (
		jobDone      = job.done
		jobError     = false
		jobErrorText = ""
	)

	if job.err != nil {
		jobError = true
		jobErrorText = job.err.Error()
	}

	return jobDone, jobError, jobErrorText, nil
}
