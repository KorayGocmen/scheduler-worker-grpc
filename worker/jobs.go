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
	done        bool
	err         error
}

func startScript(command, path string) (bool, string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	timestamp := time.Now().Format("20060102150405")
	outFilePath := fmt.Sprintf("%s.out", timestamp)
	newJob := job{
		command:     command,
		path:        path,
		outFilePath: outFilePath,
		done:        false,
		err:         nil,
	}

	jobs[path] = &newJob

	cmd := exec.Command(command, path)

	outfile, err := os.Create(outFilePath)
	if err != nil {
		jobs[path].done = true
		jobs[path].err = err
		return false, err.Error()
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	if err = cmd.Start(); err != nil {
		jobs[path].done = true
		jobs[path].err = err
		return false, err.Error()
	}

	return true, ""
}

func stopScript(path string) (bool, string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()

	job, found := jobs[path]
	if !found {
		return false, "Job not found."
	}

	job.done = true
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
