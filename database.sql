CREATE TABLE currency_pairs (
    pair_id SERIAL PRIMARY KEY,
    base_currency TEXT NOT NULL,   -- e.g., BTC
    quote_currency TEXT NOT NULL,  -- e.g., USDT
    UNIQUE (base_currency, quote_currency)
);

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


