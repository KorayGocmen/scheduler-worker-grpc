# WORKER

Workers exposes a GRPC-server to communicate with the scheduler. Main job of workers are to run specific jobs requested by the scheduler and report on those jobs. Jobs can be ruby/python/bash scripts or any executables available on the worker machine.

## Data Structures
---

### Jobs

```golang
// jobsMutex is the lock to access jobs map.
// jobs is the map that holds current/past jobs.
//    - key: job id
//    - value: pointer to the created job object.
var (
	jobsMutex = &sync.Mutex{}
	jobs      = make(map[string]*job)
)

// job holds information about the ongoing or past jobs,
// that were triggered by the scheduler.
//    - id: UUID assigned by the worker and sent back to the scheduler.
//    - command: command which the scheduler run the job with
//    - path: path to the job file/executable sent by the scheduler.
//    - outFilePath: file path to where the output of the job will be piped.
//    - cmd: pointer to the cmd.Exec command to get job status etc.
//    - done: whether if job is done (default false)
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
```

## Config
---

```toml
[grpc_server]
addr = "127.0.0.1:30000"
use_tls = false
crt_file = "server.pem"
key_file = "server.key"

[scheduler]
addr = "127.0.0.1:50000"
```

`grpc_server`:
  - `addr`: Address on which the GRPC server will be run.
  - `use_tls`: Whether the GRPC server should use TLS. If `true`, `crt_file` and `key_file` should be provided.
  - `crt_file`: Path to the certificate file for TLS.
  - `key_file`: Path to the key file for TLS.

`scheduler`:
  - `addr`: Address on which the GRPC server of the scheduler is run.
  

## GRPC Server
---

#### Methods:
```
service Worker {
  rpc StartJob(StartJobReq) returns (StartJobRes) {}
  rpc StopJob(StopJobReq) returns (StopJobRes) {}
  rpc QueryJob(QueryJobReq) returns (QueryJobRes) {}
  rpc StreamJob(StreamJobReq) returns (stream StreamJobRes) {}
}
```

### Variables:
```
message StartJobReq {
  string command = 1;
  string path = 2;
}

message StartJobRes {
  string jobID = 1;
}

message StopJobReq {
  string jobID = 1;
}

message StopJobRes {
}

message QueryJobReq {
  string jobID = 1;
}

message QueryJobRes {
  bool done = 1;
  bool error = 2;
  string errorText = 3;
}

message StreamJobReq {
  string path = 1;
}

message StreamJobRes {
  string output = 1;
}
```

`StartJob`:
- Called by the scheduler to start a new job. 
- `StartJobReq` has the `command` to run the job with and the `path` to the script/executable.
- `StartJobRes` returns a `jobID` created by the worker that later used by the scheduler to specify the created job.

`StopJob`:
- Called by the scheduler to stop a job with the given `jobID`.

`QueryJob`:
- Called by the scheduler to query a job with the given `jobID`.
- Returns if the job is `done`, or if there was an `error` and if there was an error, what the `errorText` output.

`StreamJob`:
- Called by the scheduler to get streaming output of the job with the given `jobID`.