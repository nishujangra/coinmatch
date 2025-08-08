# Sample API Usage

## 📘 API: Add New Currency Pair

### `POST /api/pairs`

Add a new trading pair (e.g., BTC/USDT) to the system. This route is **restricted to admin users** and requires a valid Bearer token.

---

## 🔐 Authentication

This endpoint requires an `Authorization` header with a **Bearer token**. The expected token must match the `ADMIN_TOKEN` set in the `.env` file.

**Header Example:**

```http
Authorization: Bearer <ADMIN_TOKEN>
```

---

## 📤 Request Body

Send a JSON object containing the base and quote currency:

```json
{
  "base": "BTC",
  "quote": "USDT"
}
```

| Field | Type   | Required | Description                     |
| ----- | ------ | -------- | ------------------------------- |
| base  | string | Yes      | The base currency (e.g., BTC)   |
| quote | string | Yes      | The quote currency (e.g., USDT) |

---

## ✅ Success Response

**Status:** `201 Created`

```json
{
  "message": "added successfully"
}
```

---

## ❌ Error Responses

### 1. Missing or Invalid Token

**Status:** `401 Unauthorized`

```json
{
  "error": "Invalid Auth Token"
}
```

---

### 2. Missing Authorization Header

**Status:** `401 Unauthorized`

```json
{
  "error": "Authentication Header missing"
}
```

---

### 3. Invalid Header Format

**Status:** `401 Unauthorized`

```json
{
  "error": "Invalid Authorization header format"
}
```

---

### 4. Invalid JSON Payload

**Status:** `400 Bad Request`

```json
{
  "error": "Invalid JSON body"
}
```

---

## 🧪 Example Requests

### ✅ Successful `curl` Request

```bash
curl -X POST http://localhost:8080/api/pairs \
  -H "Authorization: Bearer mysecrettoken" \
  -H "Content-Type: application/json" \
  -d '{"base": "BTC", "quote": "USDT"}'
```

---

### ❌ Missing Token

```bash
curl -X POST http://localhost:8080/api/pairs \
  -H "Content-Type: application/json" \
  -d '{"base": "BTC", "quote": "USDT"}'
```

---

### ❌ Invalid JSON

```bash
curl -X POST http://localhost:8080/api/pairs \
  -H "Authorization: Bearer mysecrettoken" \
  -H "Content-Type: application/json" \
  -d '{"base": "BTC"}'
```