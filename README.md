# ⚔️ Limit Order Book Matching Engine

A Go-based matching engine and REST API service for managing limit orders for cryptocurrency trading. This MVP supports placing buy/sell limit orders, viewing the order book, and order matching based on price-time priority (FIFO).

---

## 🎯 Features

### ✅ 1. Currency Management
- Admin can add new trading pairs (e.g., BTC/USDT).
- Endpoint: `POST /api/pairs`
- Request:
```json
{
  "base": "BTC",
  "quote": "USDT"
}
````

---

### ✅ 2. Place Limit Orders

* Users can place **buy** or **sell** limit orders.
* Only one pair supported for MVP (e.g., BTC/USDT).
* Endpoint: `POST /api/orders`
* Request:

```json
{
  "pair": "BTC/USDT",
  "side": "buy", 
  "price": 29500.00,
  "quantity": 0.5,
  "user_id": 123
}
```

* Orders stored in:

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

### ✅ 3. Matching Engine Logic

* Triggered on order placement.

* Buy Order: Matches with **lowest-priced** sell orders.

* Sell Order: Matches with **highest-priced** buy orders.

* Prioritized by:

  1. **Best price**
  2. **Earliest timestamp (FIFO)**

* Matching rules:

```go
if buy.Price >= sell.Price {
    matchedQty = min(buy.Quantity, sell.Quantity)
    // update both orders accordingly
}
```

* Statuses:

  * `filled`: completely matched
  * `partial`: partially matched, pending
  * `open`: unmatched
  * `cancelled`: manually removed

---

### ✅ 4. View Order Book

* Endpoint: `GET /api/orderbook?pair=BTC/USDT&depth=10`
* Response:

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

* Endpoint: `GET /api/orders?user_id=123`
* Shows user’s open, filled, and partially filled orders.

---

## 🧠 Data Structures

* **Two Priority Queues** per currency pair:

  * **Buy Queue (Max-Heap)**: Highest price first
  * **Sell Queue (Min-Heap)**: Lowest price first
* Each queue maintains FIFO by using timestamp.

---

## 🧱 Folder Structure

```plaintext
limitbook/
├── cmd/server/main.go
├── internal/
│   ├── api/               # API Handlers
│   ├── engine/            # Matching Logic
│   ├── models/            # DB Models
│   ├── db/                # DB connection & schema
│   └── utils/             # Logger, helpers
├── pkg/config/            # Config loader
├── migrations/            # SQL migrations
├── scripts/               # Setup scripts
├── api.postman_collection.json
├── go.mod
├── README.md
└── Makefile               # Optional build/test/run targets
```

---

## 🚀 How It Works

1. Admin adds trading pairs.
2. User places a buy/sell order.
3. Engine matches it with opposing side using priority rules.
4. Orders updated in DB based on match result.
5. Order book and user orders accessible via APIs.

---

## 💡 Stretch Goals (Optional)

| Feature                                   | Status |
| ----------------------------------------- | ------ |
| Cancel Order (`DELETE /api/orders/:id`)   | 🔜     |
| Multi-currency support                    | 🔜     |
| WebSocket for live order book             | 🔜     |
| Persist partial orders                    | 🔜     |
| Use container/heap in-memory for matching | ✅      |

---

## 🛠️ Setup

### Prerequisites

* Go 1.21+
* PostgreSQL 13+
* Make / Bash (optional)

### Run Locally

```bash
git clone https://github.com/yourname/limitbook.git
cd limitbook

# Set environment variables or use .env file
export DB_URL=postgres://user:pass@localhost:5432/limitbook

# Run the server
go run cmd/server/main.go
```

---

## 🧪 Testing

Use the included **Postman collection** to test:

* `/api/pairs`
* `/api/orders`
* `/api/orderbook`
* `/api/orders?user_id=...`

---

## 📌 Tech Stack

* Go (Golang)
* PostgreSQL
* REST API (net/http)
* In-memory matching engine
* Priority Queues (`container/heap`)

---

## 📄 License

MIT License. Free to use, modify, and distribute.

---

## 🙌 Author

Built with ⚡ by [Nishant](https://github.com/nishujangra)