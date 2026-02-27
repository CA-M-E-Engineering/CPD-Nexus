# IoT Bridge Communication Documentation

This document outlines the communication protocol between the unified backend and the IoT Bridge via WebSocket.

## Common Message Envelope

All communication across the bridge uses a standardized message envelope.

**Direction:** Both (Request and Response)

```json
{
  "meta": {
    "request_id": "req-20260227115854",
    "sent_at": "2026-02-27T11:58:54Z"
  },
  "action": "ACTION_NAME",
  "payload": { ... }
}
```

---

## 1. Attendance Query (`GET_ATTENDANCE`)

Used by the backend to fetch attendance records for a specific worker across a set of devices for a specific time range.

### Request (Backend -> Bridge)
- **Action:** `GET_ATTENDANCE`
- **Payload:**
```json
{
  "worker_id": "w20260225135067",
  "devices": ["SN-DEV-001", "SN-DEV-002"],
  "start_time": "2026-02-27T00:00:00Z",
  "end_time": "2026-02-27T11:58:54Z"
}
```

### Response (Bridge -> Backend)
- **Action:** `GET_ATTENDANCE_RESPONSE`
- **Payload:**
```json
{
  "code": 200,
  "msg": "Success",
  "content": {
    "worker_id": "w20260225135067",
    "records": [
      {
        "time_in": "2026-02-27T08:30:00Z",
        "time_out": "2026-02-27T17:45:00Z"
      }
    ]
  }
}
```

---

## 2. User Registration (`REGISTER_USER`)

Used to push a new worker's details (including biometrics or card data) to a list of devices.

### Request (Backend -> Bridge)
- **Action:** `REGISTER_USER`
- **Payload:**
```json
{
  "devices": ["SN-DEV-001"],
  "user": {
    "employee_no": "w20260225135067",
    "name": "John Doe",
    "user_type": "normal",
    "validity": {
      "start_time": "2026-02-27T00:00:00Z",
      "end_time": "2027-02-27T23:59:59Z"
    },
    "authentication": {
      "card": {
        "card_no": "987654321",
        "card_type": "normal"
      },
      "face": {
        "face_id": "101",
        "face_url": "https://storage.nxs.com/faces/w20260225135067.jpg"
      }
    }
  }
}
```

### Expected Response (Bridge -> Backend)
- **Action:** `REGISTER_USER_RESPONSE`
- **Payload:**
```json
{
  "code": 200,
  "msg": "User successfully registered on 1/1 devices",
  "content": null
}
```

---

## 3. User Update (`UPDATE_USER`)

Used to update an existing worker's details or credentials on their assigned devices.

### Request (Backend -> Bridge)
- **Action:** `UPDATE_USER`
- **Payload:** (Identical structure to `REGISTER_USER`)
```json
{
  "devices": ["SN-DEV-001"],
  "user": {
    "employee_no": "w20260225135067",
    "name": "John Doe Updated",
    "user_type": "normal",
    "validity": {
      "start_time": "2026-02-27T00:00:00Z",
      "end_time": "2028-02-27T23:59:59Z"
    },
    "authentication": {
      "card": { "card_no": "11223344", "card_type": "normal" },
      "face": { "face_id": "101", "face_url": "..." }
    }
  }
}
```

### Expected Response (Bridge -> Backend)
- **Action:** `UPDATE_USER_RESPONSE`
- **Payload:**
```json
{
  "code": 200,
  "msg": "User successfully updated",
  "content": null
}
```

---

## Unhappy Paths / Error Responses

In case of errors, the Bridge returns a non-200 code and a descriptive message.

### 1. Worker Not Found (404)
Occurs when the `worker_id` requested in `GET_ATTENDANCE` is not recognized by the bridge or has no history on the specified devices.

**Response Payload:**
```json
{
  "code": 404,
  "msg": "Worker w20260225135067 not found in bridge registry",
  "content": null
}
```

### 2. Device Sync Failure (500)
Occurs during `REGISTER_USER` or `UPDATE_USER` if the bridge fails to communicate with the physical hardware (e.g., device offline).

**Response Payload:**
```json
{
  "code": 500,
  "msg": "Failed to sync user to 1/1 devices. Device SN-DEV-001 is offline.",
  "content": null
}
```

---

## Summary of Communication

| Action | Description | Response Action |
| :--- | :--- | :--- |
| `GET_ATTENDANCE` | Fetches attendance logs for a worker | `GET_ATTENDANCE_RESPONSE` |
| `REGISTER_USER` | Syncs new worker to devices | `REGISTER_USER_RESPONSE` |
| `UPDATE_USER` | Updates existing worker on devices | `UPDATE_USER_RESPONSE` |

> [!NOTE]
> All timestamps follow the **RFC3339** format (e.g., `2006-01-02T15:04:05Z`). The `devices` array expects Unique Serial Numbers (SN) as identified in the `devices` table.
