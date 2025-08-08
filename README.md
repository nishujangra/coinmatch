# ⚔️ Coinmatch — Limit Order Book Matching Engine

**Coinmatch** is a Go-based limit order book and matching engine for cryptocurrency trading. It enables users to place **buy/sell limit orders**, and automatically matches them based on **price-time priority (FIFO)**. The system is built with modularity and extendability in mind, using PostgreSQL for persistence and in-memory matching queues for performance.

---

## 🎯 MVP Features

### ✅ 1. Currency Management (Admin Only)

Admins can add supported trading pairs like `BTC/USDT`.

**Endpoint:** `POST /api/pairs`  
**Request Body:**
```json
{
  "base": "BTC",
  "quote": "USDT"
}
````

---

### ✅ 2. Place Limit Orders

Users can place **limit buy/sell orders** on supported pairs.

**Endpoint:** `POST /api/orders`
**Request Body:**

```json
{
  "pair": "BTC/USDT",
  "side": "buy", 
  "price": 29500.00,
  "quantity": 0.5,
  "user_id": 123
}
```

**Order Table Schema:**

```sql
CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  user_id INT,
  pair TEXT,
  side TEXT CHECK (side IN ('buy', 'sell')),
  price DECIMAL,
  quantity DECIMAL,
  filled_quantity DECIMAL DEFAULT 0,
  status TEXT CHECK (status IN ('open', 'filled', 'partial', 'cancelled')),
  created_at TIMESTAMP DEFAULT now()
);
```

---

### ✅ 3. Matching Engine (Runs On Order Placement)

Matching is triggered immediately after an order is placed.

**Logic:**

* **Buy Orders** match with lowest-priced **Sell Orders**
* **Sell Orders** match with highest-priced **Buy Orders**
* Matches by:

  1. Best Price
  2. Earliest Time (FIFO)
* Status updates:

  * `filled`: fully matched
  * `partial`: partially matched
  * `open`: unmatched
  * `cancelled`: manually removed

**Sample Rule:**

```go
if buy.Price >= sell.Price {
    matchedQty = min(buy.Quantity, sell.Quantity)
    // update both orders accordingly
}
```

---

### ✅ 4. View Order Book

Returns top N buy/sell orders for a pair.

**Endpoint:** `GET /api/orderbook?pair=BTC/USDT&depth=10`
**Response:**

```json
{
  "buy": [
    { "price": 29400, "quantity": 1.5 },
    ...
  ],
  "sell": [
    { "price": 29600, "quantity": 2.0 },
    ...
  ]
}
```

---

### ✅ 5. Get User Orders

Retrieve all open, filled, and partial orders by user.

**Endpoint:** `GET /api/orders?user_id=123`

---

## 🧠 Matching Engine Data Structures

* Two in-memory **Priority Queues** per pair:

  * **Buy Queue:** Max-Heap (by price), FIFO (by time)
  * **Sell Queue:** Min-Heap (by price), FIFO (by time)
* Matching runs as part of the `/api/orders` handler

---

## 🧱 Project Structure

```
coinmatch/
├── cmd/server/main.go         # Entry point
├── internal/
│   ├── api/                   # HTTP Handlers
│   ├── engine/                # Matching engine logic
│   ├── models/                # DB structs & queries
│   ├── db/                    # Connection + migrations
│   └── utils/                 # Logger, helpers, errors
├── pkg/config/                # Configuration loader
├── migrations/                # SQL setup scripts
├── api.postman_collection.json
├── go.mod
├── Makefile                   # Build/test scripts
└── README.md
```

---

## 🚀 How It Works

1. Admin adds a currency pair via `/api/pairs`
2. User submits a limit order via `/api/orders`
3. Matching engine compares it against the opposite side queue
4. Orders are matched and stored in the DB
5. Users and viewers can query order book and order history

---

## 🔄 API Summary

| Endpoint                  | Method | Description                    |
| ------------------------- | ------ | ------------------------------ |
| `/api/pairs`              | POST   | Add a new trading pair (Admin) |
| `/api/orders`             | POST   | Place a new limit order        |
| `/api/orderbook?pair=...` | GET    | View current order book        |
| `/api/orders?user_id=...` | GET    | Get all orders by user ID      |
| `/api/orders/{id}`        | DELETE | (Optional) Cancel open order   |

---

## 💡 Stretch Goals

| Feature                                 | Status |
| --------------------------------------- | ------ |
| Cancel Order (`DELETE /api/orders/:id`) | 🔜     |
| Multi-currency pair support             | 🔜     |
| WebSocket order book stream             | 🔜     |
| Persistent queue snapshots              | 🔜     |
| Heap-based matching engine (Go PQ)      | ✅      |

---

## 🛠️ Setup

### Prerequisites

* Go 1.21+
* PostgreSQL 13+
* (Optional) `make`, `.env`

### Run Locally

```bash
git clone https://github.com/yourname/coinmatch.git
cd coinmatch

# Set environment variable
export DB_URL=postgres://user:pass@localhost:5432/coinmatch

# Run the server
go run cmd/server/main.go
```

---

## 🧪 API Testing

Use the provided **Postman collection**: `api.postman_collection.json`

Test:

* Add trading pairs
* Place buy/sell orders
* View live order book
* Query orders by user

---

## 📌 Tech Stack

* **Golang** — REST API & Engine
* **PostgreSQL** — Persistent storage
* **In-Memory Heaps** — Matching logic (`container/heap`)
* **net/http + chi** — Routing layer
* **Modular Design** — `internal/`, `cmd/`, `pkg/`

---

## 📄 License

MIT License — Free for personal and commercial use.

---

## 🙌 Author

Built with ⚡ by [Nishant Jangra](https://github.com/nishujangra)