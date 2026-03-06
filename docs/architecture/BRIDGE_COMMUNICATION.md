# CPD-Nexus — IoT Bridge Communication Protocol

This document defines the WebSocket message protocol used between the CPD-Nexus backend and the IoT device bridge. All messages use a standardised JSON envelope.

---

## Message Envelope

Every message — in both directions — follows this structure:

```json
{
  "meta": {
    "request_id": "req-20260301120530|w20260225135067",
    "sent_at": "2026-03-01T12:05:30Z",
    "auth_token": "client-bridge-secret-token"
  },
  "action": "ACTION_NAME",
  "payload": { ... }
}
```

| Field | Description |
|---|---|
| `meta.request_id` | Unique request identifier. May include contextual data after `\|` (e.g. worker ID) to correlate async responses. |
| `meta.sent_at` | RFC3339 timestamp of when the message was sent. |
| `meta.auth_token` | **Required** for all write operations (`REGISTER_USER`, `UPDATE_USER`). The bridge validates this token before modifying device state. |
| `action` | The operation name (see actions below). |
| `payload` | Action-specific data body. |

---

## Actions

### 1. `GET_ATTENDANCE` — Fetch Attendance Records

Sent by the backend to request attendance logs for a specific worker across a set of devices.

**Direction:** Backend → Bridge

```json
{
  "meta": { ... },
  "action": "GET_ATTENDANCE",
  "payload": {
    "worker_id": "w20260225135067",
    "devices": ["SN-DEV-001", "SN-DEV-002"],
    "start_time": "2026-03-01T00:00:00Z",
    "end_time": "2026-03-01T23:59:59Z"
  }
}
```

**Response Action:** `GET_ATTENDANCE_RESPONSE` (Bridge → Backend)

```json
{
  "meta": { "request_id": "req-20260301120530|w20260225135067", ... },
  "action": "GET_ATTENDANCE_RESPONSE",
  "payload": {
    "code": 200,
    "msg": "Success",
    "content": {
      "worker_id": "w20260225135067",
      "records": [
        {
          "time_in": "2026-03-01T08:30:00Z",
          "time_out": "2026-03-01T17:45:00Z"
        }
      ]
    }
  }
}
```

**Backend behaviour on receipt:**
- The `AttendanceHandler` parses each record and calls `AttendanceService.ProcessBridgeAttendance()`.
- A new row is inserted into `attendance` with `status = 'pending'`.
- The worker ID is extracted from the `request_id` field (after the `|` separator).

---

### 2. `REGISTER_USER` — Register New Worker on Devices

Sent to push a new worker's biometric/card credentials to one or more devices.

**Direction:** Backend → Bridge

```json
{
  "meta": { "auth_token": "client-bridge-secret-token", ... },
  "action": "REGISTER_USER",
  "payload": {
    "devices": ["SN-DEV-001"],
    "user": {
      "employee_no": "w20260225135067",
      "name": "John Doe",
      "user_type": "normal",
      "validity": {
        "start_time": "2026-03-01T00:00:00Z",
        "end_time":   "2027-03-01T23:59:59Z"
      },
      "authentication": {
        "card": {
          "card_no":   "987654321",
          "card_type": "normal"
        },
        "face": {
          "face_id":  "101",
          "face_url": "https://storage.example.com/faces/w20260225135067.jpg"
        }
      }
    }
  }
}
```

**Response Action:** `REGISTER_USER_RESPONSE` (Bridge → Backend)

```json
{
  "action": "REGISTER_USER_RESPONSE",
  "payload": {
    "code": 200,
    "msg": "User successfully registered on 1/1 devices",
    "content": null
  }
}
```

**Backend behaviour on receipt:**
- The `UserSyncResponseHandler` processes the response.
- **Only on HTTP 200**: the worker's `is_synced` flag is set to `1` (Synced).
- On any non-200 response, `is_synced` is left unchanged so the worker remains in the retry queue.

---

### 3. `UPDATE_USER` — Update Existing Worker on Devices

Sent when a registered worker's name, credentials, or validity period has changed.

**Direction:** Backend → Bridge

```json
{
  "meta": { "auth_token": "client-bridge-secret-token", ... },
  "action": "UPDATE_USER",
  "payload": {
    "devices": ["SN-DEV-001"],
    "user": {
      "employee_no": "w20260225135067",
      "name": "John Doe (Updated)",
      "user_type": "normal",
      "validity": {
        "start_time": "2026-03-01T00:00:00Z",
        "end_time":   "2028-03-01T23:59:59Z"
      },
      "authentication": {
        "card": { "card_no": "11223344", "card_type": "normal" },
        "face": { "face_id": "101", "face_url": "..." }
      }
    }
  }
}
```

**Response Action:** `UPDATE_USER_RESPONSE` (Bridge → Backend)

```json
{
  "action": "UPDATE_USER_RESPONSE",
  "payload": {
    "code": 200,
    "msg": "User successfully updated",
    "content": null
  }
}
```

**Backend behaviour on receipt:** Same as `REGISTER_USER_RESPONSE` — sets `is_synced = 1` only on HTTP 200.

---

## Error Responses

The bridge returns a non-200 `code` for known error conditions.

### Worker Not Found (404)
```json
{
  "payload": {
    "code": 404,
    "msg": "Worker w20260225135067 not found in bridge registry",
    "content": null
  }
}
```

### Device Unreachable (500)
```json
{
  "payload": {
    "code": 500,
    "msg": "Failed to sync user to 1/1 devices. Device SN-DEV-001 is offline.",
    "content": null
  }
}
```

### Unauthorized (401)
Returned when the `auth_token` in `meta` is missing or invalid. The backend does **not** update `is_synced` when this occurs — the sync is retried on the next scheduled cycle.

---

## Sync Trigger Rules

| Worker Condition | `is_synced` value | Backend Action |
|---|---|---|
| New worker with biometrics | `2` (Pending Registration) | Send `REGISTER_USER` |
| Existing worker — biometrics/name/project changed | `0` (Pending Update) | Send `UPDATE_USER` |
| Worker without biometrics/card | Any | **Not synced** — skip |
| Worker with `status = 'inactive'` | Any | **Not synced** — skip |
| Attendance status is `'submitted'` | N/A | **Not re-submitted** to SGTradeX |

---

## Summary

| Action | Direction | Response Action |
|---|---|---|
| `GET_ATTENDANCE` | Backend → Bridge | `GET_ATTENDANCE_RESPONSE` |
| `REGISTER_USER` | Backend → Bridge | `REGISTER_USER_RESPONSE` |
| `UPDATE_USER` | Backend → Bridge | `UPDATE_USER_RESPONSE` |

> [!NOTE]
> All timestamps must be **RFC3339** format (e.g. `2026-03-01T08:30:00Z`).
> The `devices` array expects device **Serial Numbers** (`sn`) as stored in the `devices` table — not internal device IDs.

> [!IMPORTANT]
> `is_synced` is only updated to `1` (Synced) when the bridge returns a `200` response code. Any other code leaves the flag unchanged so the worker is automatically retried on the next sync cycle.
