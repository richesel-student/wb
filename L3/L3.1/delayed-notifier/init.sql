CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    message TEXT,
    channel TEXT,
    recipient TEXT,
    send_at TIMESTAMP,
    status TEXT,
    retries INT,
    created_at TIMESTAMP
);