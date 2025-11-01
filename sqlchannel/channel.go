package sqlchannel

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/ali63yavari/outbox-abstraction/abstraction"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var initOnce sync.Once

type SqlEventChannel interface {
	RegisterEvent(event abstraction.OutboxEvent) error
}

type sqlEventChannel struct {
	db         *gorm.DB
	ctx        context.Context
	eventType  abstraction.OutboxEventType
	maxRetries int
	batchSize  int
	interval   time.Duration
	handler    abstraction.OutboxEventHandler
}

func NewsqlEventChannel(ctx context.Context,
	db *gorm.DB,
	eventType abstraction.OutboxEventType,
	maxRetries *int,
	batchSize *int,
	interval *time.Duration,
	handler abstraction.OutboxEventHandler,
) SqlEventChannel {
	mr := 3
	if maxRetries != nil && *maxRetries > 0 {
		mr = *maxRetries
	}

	bs := 10
	if batchSize != nil && *batchSize > 0 {
		bs = *batchSize
	}

	itr := 60 * time.Second
	if interval != nil && *interval > 0 {
		itr = *interval
	}

	initOnce.Do(func() {
		if err := db.AutoMigrate(&EventModel{}); err != nil {
			log.Printf("ERROR: sql channel migration failed: %v", err)
		}
	})

	ch := &sqlEventChannel{
		db:         db,
		ctx:        ctx,
		eventType:  eventType,
		maxRetries: mr,
		batchSize:  bs,
		interval:   itr,
		handler:    handler,
	}

	go ch.run()

	return ch
}

func (c *sqlEventChannel) RegisterEvent(event abstraction.OutboxEvent) error {
	md := FromOutboxEvent(&event)
	err := c.db.Model(EventModel{}).Create(&md).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *sqlEventChannel) run() {
	timer := time.NewTicker(c.interval)
	defer timer.Stop()

	for {
		select {
		case <-c.ctx.Done():
			log.Printf("[%s] Outbox event channel closed.\n", c.eventType)
			return
		case <-timer.C:
			if err := c.processBatch(); err != nil {
				log.Printf("[%s] ERROR processing batch: %v\n", c.eventType, err)
			}
			log.Printf("[%s] Outbox event channel interval passed.\n", c.eventType)
		}
	}
}

func (c *sqlEventChannel) processBatch() error {
	// Fetch events in a transaction with row locking
	var events []EventModel
	err := c.db.WithContext(c.ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(EventModel{}).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("event_status = ? AND event_type = ?",
				string(abstraction.OutboxEventStatusPending), c.eventType.GetName()).
			Order("id ASC").Limit(c.batchSize).Find(&events).Error

		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}
	log.Printf("processing %d outbox events with type %s", len(events), c.eventType.GetName())

	// Process each event in its own transaction
	for i := 0; i < len(events); i++ {
		evt := &events[i]
		if err := c.handleEvent(c.ctx, evt); err != nil {
			log.Printf("event %d failed: %v", i, err)
		}
	}

	return nil
}

func (c *sqlEventChannel) handleEvent(ctx context.Context, evt *EventModel) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := c.handler(ctx, evt.ToOutboxEvent())
		now := time.Now()

		if err == nil {
			return tx.Model(&evt).Updates(map[string]interface{}{
				"event_status":  string(abstraction.OutboxEventStatusClosed),
				"last_error":    "",
				"updated_at":    now,
				"attempt_count": evt.AttemptCount + 1,
			}).Error
		}

		// Exponential backoff
		evt.AttemptCount++
		backoff := time.Duration(math.Min(60, math.Pow(2, float64(evt.AttemptCount)))) * time.Second
		nextAttempt := now.Add(backoff)

		update := map[string]interface{}{
			"event_status":    string(abstraction.OutboxEventStatusPending),
			"last_error":      err.Error(),
			"attempt_count":   evt.AttemptCount,
			"next_attempt_at": nextAttempt,
			"updated_at":      now,
		}
		if evt.AttemptCount >= c.maxRetries {
			update["event_status"] = string(abstraction.OutboxEventStatusFailed)
		}

		return tx.Model(&evt).Updates(update).Error
	})
}
