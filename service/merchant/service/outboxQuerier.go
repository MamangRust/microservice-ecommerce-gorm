package service

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// outboxQuerier adapts GORM to the outbox.OutboxQuerier contract, replacing the
// sqlc-generated outbox queries. Claiming uses FOR UPDATE SKIP LOCKED so
// concurrent relay instances never publish the same event twice.
type outboxQuerier struct {
	db *gorm.DB
}

func NewOutboxQuerier(db *gorm.DB) outbox.OutboxQuerier {
	if db == nil {
		return nil
	}
	return &outboxQuerier{db: db}
}

func (q *outboxQuerier) CreateOutboxEvent(ctx context.Context, params outbox.CreateOutboxEventParams) (outbox.OutboxEvent, error) {
	now := time.Now()
	event := models.OutboxEvent{
		Topic:         params.Topic,
		EventKey:      params.EventKey,
		Payload:       params.Payload,
		Status:        "pending",
		NextAttemptAt: now,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	if err := q.db.WithContext(ctx).Create(&event).Error; err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outboxEventToDomain(event), nil
}

func (q *outboxQuerier) ClaimPendingOutboxEvents(ctx context.Context, params outbox.ClaimPendingOutboxEventsParams) ([]outbox.OutboxEvent, error) {
	// Claim atomically: take the lease (next_attempt_at) with SKIP LOCKED so
	// concurrent relays never claim the same event.
	err := q.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND next_attempt_at <= ?", "pending", time.Now()).
		Order("outbox_id").
		Limit(int(params.Limit)).
		Update("next_attempt_at", params.NextAttemptAt).Error
	if err != nil {
		return nil, err
	}

	// Re-query the claimed events (those holding the lease).
	var events []models.OutboxEvent
	err = q.db.WithContext(ctx).
		Where("status = ? AND next_attempt_at = ? AND next_attempt_at > ?", "pending", params.NextAttemptAt, time.Now()).
		Order("outbox_id").
		Find(&events).Error
	if err != nil {
		return nil, err
	}

	result := make([]outbox.OutboxEvent, 0, len(events))
	for _, e := range events {
		result = append(result, outboxEventToDomain(e))
	}
	return result, nil
}

func (q *outboxQuerier) MarkOutboxEventFailed(ctx context.Context, params outbox.MarkOutboxEventFailedParams) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).
		Where("outbox_id = ? AND status = ?", params.OutboxID, "pending").
		First(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}

	event.Attempts++
	event.NextAttemptAt = params.NextAttemptAt
	now := time.Now()
	event.UpdatedAt = &now

	err = q.db.WithContext(ctx).Save(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outboxEventToDomain(event), nil
}

func (q *outboxQuerier) MarkOutboxEventDelivered(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).
		Where("outbox_id = ? AND status = ?", outboxID, "pending").
		First(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}

	event.Status = "delivered"
	now := time.Now()
	event.UpdatedAt = &now

	err = q.db.WithContext(ctx).Save(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outboxEventToDomain(event), nil
}

func (q *outboxQuerier) MarkOutboxEventDead(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).
		Where("outbox_id = ? AND status = ?", outboxID, "pending").
		First(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}

	event.Status = "dead"
	now := time.Now()
	event.UpdatedAt = &now

	err = q.db.WithContext(ctx).Save(&event).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outboxEventToDomain(event), nil
}

func (q *outboxQuerier) DeleteOldOutboxEvents(ctx context.Context, cutoff time.Time) (int64, error) {
	result := q.db.WithContext(ctx).
		Where("status IN (?, ?) AND updated_at < ?", "delivered", "dead", cutoff).
		Delete(&models.OutboxEvent{})
	return result.RowsAffected, result.Error
}

func outboxEventToDomain(e models.OutboxEvent) outbox.OutboxEvent {
	var createdAt, updatedAt time.Time
	if e.CreatedAt != nil {
		createdAt = *e.CreatedAt
	}
	if e.UpdatedAt != nil {
		updatedAt = *e.UpdatedAt
	}
	return outbox.OutboxEvent{
		OutboxID:      e.OutboxID,
		Topic:         e.Topic,
		EventKey:      e.EventKey,
		Payload:       e.Payload,
		Status:        e.Status,
		Attempts:      e.Attempts,
		NextAttemptAt: e.NextAttemptAt,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}