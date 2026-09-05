package inbox

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
func Reserve(ctx context.Context, db *gorm.DB, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if db == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}

	type reserveResult struct {
		Reserved           bool
		Processed          bool
		ReservationVersion int64
	}

	// Try to INSERT with ON CONFLICT (upsert).
	// First, check if already processed.
	var processed bool
	err := db.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1 FROM consumer_inbox
			WHERE consumer_name = ? AND event_key = ? AND status = 'processed'
		)`, consumerName, eventKey).Scan(&processed).Error
	if err != nil {
		return false, false, 0, err
	}
	if processed {
		return false, true, 0, nil
	}

	// Upsert: INSERT ... ON CONFLICT DO UPDATE WHERE status <> 'processed' AND lease_until <= now()
	now := time.Now()
	leaseUntil := now.Add(1 * time.Minute)

	result := db.WithContext(ctx).Exec(`
		INSERT INTO consumer_inbox (
			consumer_name, event_key, topic, partition_id, message_offset,
			status, attempts, reservation_version, lease_until, last_error, processed_at
		) VALUES (?, ?, ?, ?, ?, 'processing', 1, 1, ?, '', NULL)
		ON CONFLICT (consumer_name, event_key) DO UPDATE
		SET status = 'processing',
			attempts = consumer_inbox.attempts + 1,
			reservation_version = consumer_inbox.reservation_version + 1,
			lease_until = ?,
			last_error = '',
			topic = EXCLUDED.topic,
			partition_id = EXCLUDED.partition_id,
			message_offset = EXCLUDED.message_offset
		WHERE consumer_inbox.status <> 'processed'
		  AND consumer_inbox.lease_until <= ?
	`, consumerName, eventKey, topic, partition, offset, leaseUntil, leaseUntil, now)

	if result.Error != nil {
		return false, false, 0, result.Error
	}

	if result.RowsAffected == 0 {
		// Already processed or still being processed by another consumer.
		return false, processed, 0, nil
	}

	// Get the reservation version.
	var reservationVersion int64
	err = db.WithContext(ctx).Raw(`
		SELECT reservation_version FROM consumer_inbox
		WHERE consumer_name = ? AND event_key = ?
	`, consumerName, eventKey).Scan(&reservationVersion).Error
	if err != nil {
		return false, false, 0, err
	}

	return true, false, reservationVersion, nil
}

func MarkProcessed(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}

	return db.WithContext(ctx).Exec(`
		UPDATE consumer_inbox
		SET status = 'processed', processed_at = current_timestamp,
			lease_until = current_timestamp, last_error = ''
		WHERE consumer_name = ? AND event_key = ?
		  AND status = 'processing' AND reservation_version = ?
	`, consumerName, eventKey, reservationVersion).Error
}

func Release(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}

	return db.WithContext(ctx).Exec(`
		UPDATE consumer_inbox
		SET status = 'pending', lease_until = current_timestamp,
			last_error = ?
		WHERE consumer_name = ? AND event_key = ?
		  AND status = 'processing' AND reservation_version = ?
	`, lastError, consumerName, eventKey, reservationVersion).Error
}

// PostgresInbox adapts a GORM DB to the outbox.ConsumerInbox contract.
// Reservation and completion are committed independently because an external
// side effect cannot share a PostgreSQL transaction with the Kafka consumer.
type PostgresInbox struct {
	db *gorm.DB
}

func NewPostgresInbox(db *gorm.DB) (*PostgresInbox, error) {
	if db == nil {
		return nil, errors.New("inbox db is nil")
	}
	return &PostgresInbox{db: db}, nil
}

func (i *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	return Reserve(ctx, i.db, consumerName, eventKey, topic, partition, offset)
}

func (i *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return MarkProcessed(ctx, i.db, consumerName, eventKey, reservationVersion)
}

func (i *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return Release(ctx, i.db, consumerName, eventKey, reservationVersion, processingErr)
}
