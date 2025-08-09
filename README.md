# ⚔️ Coinmatch — Limit Order Book Matching Engine

**Coinmatch** is a Go-based limit order book and matching engine task.

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
```

> Approach to check the user admin or not is via ADMIN_TOKEN present in `.env`

Sample request from admin:

```sh
curl -X POST /api/pairs
-H "Content-Type: application/json"
-H "Authorization: Bearer admin-secret-token"
-d '{
  "base": "BTC",
  "quote": "USDT"
}'
```

It will save to the `currency_pair` table in the database. So, it will be easy to extend later on.

Currency pair schema,
```sql
CREATE TABLE currency_pairs (
    pair_id SERIAL PRIMARY KEY,
    base_currency TEXT NOT NULL,   -- e.g., BTC
    quote_currency TEXT NOT NULL,  -- e.g., USDT
    UNIQUE (base_currency, quote_currency)
);

```

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
├── cmd/main.go         # Entry point
├── lib/
│   ├── controllers/                   # HTTP Controllers
│   ├── engine/                # Matching engine logic
│   ├── models/                # DB structs & queries
│   ├── config/                    # Connection + migrations
│   └── routes/                 # HTTP Routes
├── config/                # Configuration loader
├── database.sql                # sql tables
├── tests/postman/api.postman_collection.json # postmane collection
├── go.mod
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
| `/api/orders/{id}`        | DELETE | Cancel open order              |

---

## 🛠️ Setup

### Prerequisites

* Go 1.21+
* PostgreSQL 13+

### Run Locally

```sh
git clone https://github.com/nishuajangra/coinmatch.git
cd coinmatch

# Setup the configurations and update config.json according to your database
cp config/config.example.json config.json

# Setup environment variable and update .env after next 2 steps
touch .env
cp .env.example .env

# Run the server
go run cmd/main.go
```

---

## 🧪 API Testing

### 📦 **Postman Collection**
Complete testing suite available in `tests/postman/`:
- **Collection**: `coinmatch-api.postman_collection.json`
- **Environment**: `coinmatch-local.postman_environment.json`

### 🚀 **Quick API Examples**

1. **Add currency pair (Admin)**:

```sh
curl -X POST http://localhost:8080/api/pairs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer admin-secret-token-change-this" \
  -d '{"base": "BTC", "quote": "USDT"}'
```

2. **Place a buy order**:

```sh
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "pair": "BTC/USDT",
    "side": "buy",
    "price": 29500.00,
    "quantity": 0.5,
    "user_id": 123
  }'
```

3. **View order book**:

```sh
curl -X GET "http://localhost:8080/api/orderbook?pair=BTC/USDT&depth=10"
```

4. **Get user orders**:

```sh
curl -X GET "http://localhost:8080/api/orders?user_id=123"
```

---

## 📌 Tech Stack

* **Golang** — REST API & Engine
* **PostgreSQL** — Database
* **In-Memory Heaps** — Matching logic (`container/heap`)
* **Gin Framework** — HTTP routing and middleware
* **Modular Design** — `lib/`, `cmd/`

---

## 📄 License

[MIT License](LICENSE.md) — Free for personal and commercial use.

---

## 🙌 Author

Built with ⚡ by [Nishant](https://github.com/nishujangra)