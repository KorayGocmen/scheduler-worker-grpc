# SCHEDULER

Scheduler exposes an HTTP-server that gives a gateway to the scheduler-worker-cluster. HTTP requests to the scheduler api is translated and proxied into workers.

- Scheduler runs an HTTP server as well as a GRPC server.
- Scheduler keeps record of the available workers in the cluster.
- Workers can register and deregister on the scheduler.


## Data Structures
---

### Workers

```golang
var (
	workersMutex = &sync.Mutex{}
	workers      = make(map[string]*worker)
)

// worker holds the information about registered workers
//    - id: uuid assigned when the worker first register.
//    - addr: workers network address, later used to create grpc client to the worker
type worker struct {
	id   string
	addr string
}
```

## Config
---

```toml
[grpc_server]
addr = "127.0.0.1:50000"
use_tls = false
crt_file = "server.pem"
key_file = "server.key"

[http_server]
addr = "127.0.0.1:3000"
```

`grpc_server`:
  - `addr`: Address on which the GRPC server will be run.
  - `use_tls`: Whether the GRPC server should use TLS. If `true`, `crt_file` and `key_file` should be provided.
  - `crt_file`: Path to the certificate file for TLS.
  - `key_file`: Path to the key file for TLS.

`http_server`:
  - `addr`: Address on which the HTTP server will be run on.
  


## HTTP Server
---

Detailed API documentation for available methods can be found in [API Docs](api.md).


## GRPC Server
---

#### Methods:
```
service Scheduler {
  rpc RegisterWorker(RegisterReq) returns (RegisterRes) {}
  rpc DeregisterWorker(DeregisterReq) returns (DeregisterRes) {}
}
```

### Variables:
```
message RegisterReq {
  string address = 1;
}

message RegisterRes {
  bool success = 1;
  string workerID = 2;
}

message DeregisterReq {
  string workerID = 1;
}

message DeregisterRes {
  bool success = 1;
}
```

`RegisterWorker`:
- Called by the worker when first coming online. 
- `RegisterReq` has the address parameter which is sent by the worker. It is the address where the worker is running a GRPC server.
- `worker` instance is created for every worker that registers which contains an assigned `workerID` and the `address`.

`DeregisterWorker`:
- Called by the worker when going offline or shutting down.
- Associated worker with the `workerID` is removed from the known workers on the scheduler.

