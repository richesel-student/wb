package repository

import (
	"context"
	"database/sql"
	"time"

	"eventbooker/internal/errs"
	"eventbooker/internal/model"
)

type BookingRepository interface {
	CreateEvent(ctx context.Context, e model.Event) (int, error)
	GetEvent(ctx context.Context, id int) (model.Event, error)
	ListEvents(ctx context.Context) ([]model.Event, error)

	CreateBooking(ctx context.Context, eventID int, userID string, ttl time.Duration) (int, error)
	ConfirmBooking(ctx context.Context, bookingID int) error
	CancelBooking(ctx context.Context, bookingID int) error
	ExpireBookings(ctx context.Context, limit int) error

	FindUserBooking(ctx context.Context, eventID int, userID string) (int, error) // 👈 ДОБАВИТЬ
}

type repo struct {
	db *sql.DB
}

func New(db *sql.DB) BookingRepository {
	return &repo{db: db}
}

//////////////////////////////////////////////////////
// TX HELPER
//////////////////////////////////////////////////////

func (r *repo) withTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

//////////////////////////////////////////////////////
// EVENTS
//////////////////////////////////////////////////////

func (r *repo) CreateEvent(ctx context.Context, e model.Event) (int, error) {
	var id int

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO events (name, date, capacity, available_seats, booking_ttl)
		VALUES ($1, $2, $3, $3, $4 * interval '1 minute')
		RETURNING id
	`,
		e.Name,
		e.Date,
		e.Capacity,
		int(e.BookingTTL.Minutes()),
	).Scan(&id)

	return id, err
}

func (r *repo) GetEvent(ctx context.Context, id int) (model.Event, error) {
	var e model.Event
	var ttlSeconds float64

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, date, capacity, available_seats,
		       EXTRACT(EPOCH FROM booking_ttl)
		FROM events
		WHERE id = $1
	`, id).Scan(
		&e.ID,
		&e.Name,
		&e.Date,
		&e.Capacity,
		&e.AvailableSeats,
		&ttlSeconds,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return e, errs.ErrNotFound
		}
		return e, err
	}

	e.BookingTTL = time.Duration(ttlSeconds) * time.Second

	return e, nil
}

func (r *repo) ListEvents(ctx context.Context) ([]model.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, date, capacity, available_seats,
		        EXTRACT(EPOCH FROM booking_ttl)
		 FROM events ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []model.Event

	for rows.Next() {
		var e model.Event
		var ttl float64

		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Date,
			&e.Capacity,
			&e.AvailableSeats,
			&ttl,
		); err != nil {
			return nil, err
		}

		e.BookingTTL = time.Duration(ttl) * time.Second

		res = append(res, e)
	}

	return res, nil
}

//////////////////////////////////////////////////////
// BOOKINGS
//////////////////////////////////////////////////////

func (r *repo) CreateBooking(ctx context.Context, eventID int, userID string, ttl time.Duration) (int, error) {
	var bookingID int

	err := r.withTx(ctx, func(tx *sql.Tx) error {
		var seats int

		// блокируем event
		if err := tx.QueryRowContext(ctx,
			`SELECT available_seats FROM events WHERE id=$1 FOR UPDATE`,
			eventID,
		).Scan(&seats); err != nil {
			return err
		}

		if seats <= 0 {
			return errs.ErrNoSeats
		}

		// уменьшаем seats
		res, err := tx.ExecContext(ctx,
			`UPDATE events SET available_seats = available_seats - 1 WHERE id=$1`,
			eventID,
		)
		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil || rows == 0 {
			return errs.ErrNotFound
		}

		expiresAt := time.Now().UTC().Add(ttl)

		return tx.QueryRowContext(ctx,
			`INSERT INTO bookings(event_id,user_id,status,created_at,expires_at)
			 VALUES ($1,$2,'pending',NOW(),$3)
			 RETURNING id`,
			eventID, userID, expiresAt,
		).Scan(&bookingID)
	})

	return bookingID, err
}

//////////////////////////////////////////////////////

func (r *repo) ConfirmBooking(ctx context.Context, bookingID int) error {
	return r.withTx(ctx, func(tx *sql.Tx) error {
		var status string
		var expiresAt time.Time

		err := tx.QueryRowContext(ctx,
			`SELECT status, expires_at FROM bookings WHERE id=$1 FOR UPDATE`,
			bookingID,
		).Scan(&status, &expiresAt)

		if err == sql.ErrNoRows {
			return errs.ErrNotFound
		}
		if err != nil {
			return err
		}

		switch status {
		case "confirmed":
			return errs.ErrAlreadyConfirmed
		case "canceled":
			return errs.ErrCanceled
		case "pending":
			if expiresAt.Before(time.Now()) {
				return errs.ErrExpired
			}
		}

		_, err = tx.ExecContext(ctx,
			`UPDATE bookings SET status='confirmed' WHERE id=$1`,
			bookingID,
		)
		return err
	})
}

//////////////////////////////////////////////////////

func (r *repo) CancelBooking(ctx context.Context, bookingID int) error {
	return r.withTx(ctx, func(tx *sql.Tx) error {
		var eventID int

		err := tx.QueryRowContext(ctx,
			`SELECT event_id FROM bookings
			 WHERE id=$1 AND status='pending'
			 FOR UPDATE`,
			bookingID,
		).Scan(&eventID)

		if err == sql.ErrNoRows {
			return errs.ErrNotFound
		}
		if err != nil {
			return err
		}

		// cancel booking
		_, err = tx.ExecContext(ctx,
			`UPDATE bookings SET status='canceled' WHERE id=$1`,
			bookingID,
		)
		if err != nil {
			return err
		}

		// вернуть место (защита от overflow)
		_, err = tx.ExecContext(ctx,
			`UPDATE events
			 SET available_seats = available_seats + 1
			 WHERE id=$1 AND available_seats < capacity`,
			eventID,
		)

		return err
	})
}

//////////////////////////////////////////////////////

func (r *repo) ExpireBookings(ctx context.Context, limit int) error {
	_, err := r.db.ExecContext(ctx,
		`WITH expired AS (
			SELECT id, event_id
			FROM bookings
			WHERE status='pending' AND expires_at < NOW()
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		),
		locked_events AS (
			SELECT e.id
			FROM events e
			WHERE e.id IN (SELECT event_id FROM expired)
			ORDER BY e.id
			FOR UPDATE
		),
		updated AS (
			UPDATE bookings b
			SET status='canceled'
			FROM expired e
			WHERE b.id = e.id
			RETURNING e.event_id
		)
		UPDATE events e
		SET available_seats = available_seats + sub.cnt
		FROM (
			SELECT event_id, COUNT(*) cnt
			FROM updated
			GROUP BY event_id
		) sub
		WHERE e.id = sub.event_id`,
		limit,
	)

	return err
}

func (r *repo) FindUserBooking(ctx context.Context, eventID int, userID string) (int, error) {
	var id int

	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM bookings
		 WHERE event_id=$1 AND user_id=$2 AND status='pending'
		 ORDER BY created_at DESC
		 LIMIT 1`,
		eventID, userID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, errs.ErrNotFound
	}

	return id, err
}
