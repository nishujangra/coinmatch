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