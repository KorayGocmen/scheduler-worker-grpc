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