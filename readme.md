# Scheduler-Worker Cluster Architecture

Read full blog post here:
http://www.koraygocmen.com/single-scheduler-multiple-worker-architecture-with-grpc-and-go-part-1

# Goal

*   Build a scheduler and worker architecture that supports HTTP requests on the scheduler and translates them to GRPC requests on workers.
*   Any application outside of the GRPC cluster can make HTTP requests to trigger jobs on the cluster and can also create distributed jobs as well.


# General Overview

*   Support single scheduler and multiple workers.
*   The scheduler must be aware of all workers' states. (requires the workers to register and deregister when going offline)
*   The scheduler and the workers authenticate with each other via SSL.
*   The scheduler and workers, both need to have GRPC-APIs. 
    *   The scheduler GRPC-API support:
        *   Register a worker
        *   Deregister a worker
    *   The worker GRPC-API support:
        *   Start a job
        *   Stop a job
        *   Return the status of a job
        *   Return a stream of output for a running job.
*   The scheduler also has an HTTP-API (which will be translated and proxied to GRPC requests):
    *   Start a job on a specific worker (with worker ID, command and path to the job)
    *   Stop a job on a specific worker (with worker ID and job ID)
    *   Query a job on a specific worker (with worker ID and job ID)
    *   Return the output stream of a job a specific worker (with worker ID and job ID)


# Scheduler Overview

* [Scheduler Details](scheduler/readme.md)
* [API Docs](scheduler/api.md)
* All config parameters are specified in the config.toml file
* Support a GRPC-API (Appendix A) and an HTTP-API.
* When a worker registers, a UUID is assigned to the worker and worker details are kept in a map.
* HTTP requests are translated into GRPC requests on the specified server.
* Starting a job on a specific worker
  * Request: 
    ```
    {
      "worker_id": "71382ed1-471d-4ae3-b572-f67d178f04e9",
      "path": "worker/scripts/count.sh",
      "command": "bash"
    }
    ```
  * Response: 
    ```
    {
      "job_id": "6c26a00e-3017-11e9-b210-d663bd873d93"
    }
    ```


# Worker Overview

* [Worker Details](worker/readme.md)
* All config parameters are specified in the config.toml file
* Support a GRPC-API (Appendix B)
* Jobs are basically scripts that are held in the specified folder in the config.
* When a job is started by the scheduler, a job object is created by the worker
  * This object specifies where the output of the job will be piped.
  * Also holds which command and which path was requested.

## TODO:
---
- Finish streaming output of a job
- Change the insecure dialing for grpc
