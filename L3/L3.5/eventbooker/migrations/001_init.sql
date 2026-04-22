CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0),
    available_seats INT NOT NULL CHECK (available_seats >= 0),
    booking_ttl INTERVAL NOT NULL
);

CREATE TABLE bookings (
    id SERIAL PRIMARY KEY,
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);



CREATE INDEX idx_bookings_event_id ON bookings(event_id);
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
CREATE INDEX idx_bookings_status ON bookings(status);



CREATE UNIQUE INDEX idx_one_pending_per_user_per_event
ON bookings(event_id, user_id)
WHERE status = 'pending';