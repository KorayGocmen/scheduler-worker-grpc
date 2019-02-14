# HTTP API
---

# Start Job on Worker

#### **URL**: `/start`

#### **METHOD**: `POST`

#### **REQUEST**:

```json
{
  "worker_id": "d5feef60-3029-11e9-b210-d663bd873d93",
  "command": "bash",
  "path": "worker/scripts/count.sh"
}
```

---

#### **SUCCESS RESPONSE**:

##### **CODE**: `201 Created`

```json
{
  "job_id": "4c5ced1c-5ea9-40f8-90ce-63d09cea26f6"
}
```

---

#### **ERROR RESPONSES**:

##### **CODE**: `500 Internal Server Error`

```json
{
  "error": "worker not found"
}
```

---
---

# Stop Job on Worker

#### **URL**: `/stop`

#### **METHOD**: `POST`

#### **REQUEST**:

```json
{
  "worker_id": "d5feef60-3029-11e9-b210-d663bd873d93",
  "job_id": "4c5ced1c-5ea9-40f8-90ce-63d09cea26f6"
}
```

---

#### **SUCCESS RESPONSE**:

##### **CODE**: `200 OK`

```json
{
  "success": true
}
```

---

#### **ERROR RESPONSES**:

##### **CODE**: `500 Internal Server Error`

```json
{
  "error": "worker not found"
}
```

---
---

# Query Job on Worker

#### **URL**: `/query`

#### **METHOD**: `POST`

#### **REQUEST**:

```json
{
  "worker_id": "d5feef60-3029-11e9-b210-d663bd873d93",
  "job_id": "4c5ced1c-5ea9-40f8-90ce-63d09cea26f6"
}
```

---

#### **SUCCESS RESPONSE**:

##### **CODE**: `200 OK`

```json
{
  "done": true,
  "error": true,
  "error_text": "signal: killed"
}
```

---

#### **ERROR RESPONSES**:

##### **CODE**: `500 Internal Server Error`

```json
{
  "error": "worker not found"
}
```

---
---


## API Data Structures
---

```golang
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
	Done      bool   `json:"done"`
	Error     bool   `json:"error"`
	ErrorText string `json:"error_text"`
}

// apiError is used as a generic api response error
type apiError struct {
	Error string `json:"error"`
}

```