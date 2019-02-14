package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	jobsMutex = &sync.Mutex{}
	jobs      = make(map[string]*job)
)

type job struct {
	command     string
	path        string
	outFilePath string
	cmd         *exec.Cmd
	done        bool
	err         error
}

func startScript(command, path string) (bool, string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	timestamp := time.Now().Format("20060102150405")
	outFilePath := fmt.Sprintf("%s.out", timestamp)

	outfile, err := os.Create(outFilePath)
	if err != nil {
		return false, err.Error()
	}
	defer outfile.Close()

	cmd := exec.Command(command, path)
	cmd.Stdout = outfile

	if err = cmd.Start(); err != nil {
		return false, err.Error()
	}

	newJob := job{
		command:     command,
		path:        path,
		outFilePath: outFilePath,
		cmd:         cmd,
		done:        false,
		err:         nil,
	}
	jobs[path] = &newJob

	// Get the status of the job async.
	go func() {
		if err := cmd.Wait(); err != nil {
			newJob.err = err
		}
		newJob.done = true
	}()

	return true, ""
}

func stopScript(path string) (bool, string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, found := jobs[path]
	if !found {
		return false, "Job not found."
	}

	if err := job.cmd.Process.Kill(); err != nil {
		return false, err.Error()
	}

	return true, ""
}

func queryScript(path string) (bool, string, bool) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, found := jobs[path]
	if !found {
		return false, "Job not found.", false
	}

	return true, "", job.done
}
