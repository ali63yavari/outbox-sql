package pgsqlchannel

import (
	"time"

	"github.com/arash/outbox_abstraction/abstraction"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type EventModel struct {
	ID            uuid.UUID   `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	EventID       string      `gorm:"not null"`
	EventType     string      `gorm:"index"`
	EventStatus   string      `gorm:"index"`
	AggregateType string      `gorm:"index"`
	AggregateID   string      `gorm:"index"`
	Payload       interface{} `gorm:"serializer:json;type:jsonb"`

	AttemptCount  int
	LastError     string
	NextAttemptAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func FromOutboxEvent(event *abstraction.OutboxEvent) EventModel {
	return EventModel{
		EventID:       event.EventID,
		EventStatus:   event.EventStatus,
		EventType:     event.EventType,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,

		Payload:       event.Payload,
		AttemptCount:  0,
		NextAttemptAt: time.Now(),
		LastError:     "",

		CreatedAt: event.CreatedAt,
	}
}

func (em *EventModel) ToOutboxEvent() *abstraction.OutboxEvent {
	//TODO: maybe some isuuse to convert payload
	return &abstraction.OutboxEvent{
		EventID:       em.EventID,
		EventStatus:   em.EventStatus,
		EventType:     em.EventType,
		AggregateType: em.AggregateType,
		AggregateID:   em.AggregateID,

		Payload: em.Payload,

		CreatedAt: em.CreatedAt,
	}
}
