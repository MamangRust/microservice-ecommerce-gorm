package repository

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type outboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) OutboxRepository {
	return &outboxRepository{db: db}
}

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*models.OutboxEvent, error) {
	now := time.Now()
	event := &models.OutboxEvent{
		Topic: topic, EventKey: key, Payload: payload,
		Status: "pending", NextAttemptAt: now, CreatedAt: &now, UpdatedAt: &now,
	}
	err := r.db.WithContext(ctx).Create(event).Error
	return event, err
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*models.OutboxEvent, error) {
	var events []*models.OutboxEvent
	err := r.db.WithContext(ctx).Where("status = ? AND next_attempt_at <= ?", "pending", time.Now()).
		Order("outbox_id").Limit(limit).Find(&events).Error
	return events, err
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*models.OutboxEvent, error) {
	var events []*models.OutboxEvent
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status = ? AND next_attempt_at <= ?", "pending", time.Now()).
		Order("outbox_id").Limit(limit).Find(&events).Error
	if err != nil {
		return nil, err
	}
	// Update claimed events
	for _, e := range events {
		e.NextAttemptAt = leaseUntil
		now := time.Now()
		e.UpdatedAt = &now
	}
	if len(events) > 0 {
		r.db.WithContext(ctx).Save(&events)
	}
	return events, nil
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*models.OutboxEvent, error) {
	var event models.OutboxEvent
	err := r.db.WithContext(ctx).Where("outbox_id = ? AND status = ?", outboxID, "pending").First(&event).Error
	if err != nil { return nil, err }
	event.Status = "delivered"
	now := time.Now()
	event.UpdatedAt = &now
	err = r.db.WithContext(ctx).Save(&event).Error
	return &event, err
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*models.OutboxEvent, error) {
	var event models.OutboxEvent
	err := r.db.WithContext(ctx).Where("outbox_id = ? AND status = ?", outboxID, "pending").First(&event).Error
	if err != nil { return nil, err }
	event.Attempts++
	event.NextAttemptAt = nextAttemptAt
	now := time.Now()
	event.UpdatedAt = &now
	err = r.db.WithContext(ctx).Save(&event).Error
	return &event, err
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*models.OutboxEvent, error) {
	var event models.OutboxEvent
	err := r.db.WithContext(ctx).Where("outbox_id = ? AND status = ?", outboxID, "pending").First(&event).Error
	if err != nil { return nil, err }
	event.Status = "dead"
	now := time.Now()
	event.UpdatedAt = &now
	err = r.db.WithContext(ctx).Save(&event).Error
	return &event, err
}

func (r *outboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	result := r.db.WithContext(ctx).Where("status IN (?, ?) AND updated_at < ?", "delivered", "dead", cutoff).Delete(&models.OutboxEvent{})
	return result.RowsAffected, result.Error
}
