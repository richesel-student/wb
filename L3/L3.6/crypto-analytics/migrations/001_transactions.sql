CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    asset_id INT,
    type TEXT,
    amount_usd NUMERIC,
    price NUMERIC,
    quantity NUMERIC,
    fee NUMERIC,
    timestamp TIMESTAMP DEFAULT NOW()
);