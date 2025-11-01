package sqlchannel

import (
	"testing"
	"time"

	"github.com/arash/outbox_abstraction/abstraction"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Mock event type
type mockEventType struct {
	name string
	id   int
}

func (m mockEventType) GetName() string {
	return m.name
}

func (m mockEventType) GetID() int {
	return m.id
}

// TestEventModel for SQLite (without PostgreSQL-specific features)
type TestEventModel struct {
	ID            string `gorm:"primaryKey"`
	EventID       string `gorm:"not null"`
	EventType     string `gorm:"index"`
	EventStatus   string `gorm:"index"`
	AggregateType string `gorm:"index"`
	AggregateID   string `gorm:"index"`
	Payload       string `gorm:"type:text"`

	AttemptCount  int
	LastError     string
	NextAttemptAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TestEventModel) TableName() string {
	return "event_models"
}

// Setup test database
func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations using test model
	err = db.AutoMigrate(&TestEventModel{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Override the table name for EventModel to use the same table
	db.Table("event_models")

	return db
}

func TestEventModel_FromOutboxEvent(t *testing.T) {
	outboxEvent := &abstraction.OutboxEvent{
		EventID:       uuid.New().String(),
		EventType:     "UserCreated",
		EventStatus:   string(abstraction.OutboxEventStatusPending),
		AggregateType: "User",
		AggregateID:   "user-123",
		Payload: map[string]interface{}{
			"email": "test@example.com",
		},
		CreatedAt: time.Now(),
	}

	model := FromOutboxEvent(outboxEvent)

	// Verify all fields
	if model.EventID != outboxEvent.EventID {
		t.Errorf("EventID mismatch: expected %s, got %s", outboxEvent.EventID, model.EventID)
	}
	if model.EventType != outboxEvent.EventType {
		t.Errorf("EventType mismatch: expected %s, got %s", outboxEvent.EventType, model.EventType)
	}
	if model.EventStatus != outboxEvent.EventStatus {
		t.Errorf("EventStatus mismatch: expected %s, got %s", outboxEvent.EventStatus, model.EventStatus)
	}
	if model.AggregateType != outboxEvent.AggregateType {
		t.Errorf("AggregateType mismatch: expected %s, got %s",
			outboxEvent.AggregateType, model.AggregateType)
	}
	if model.AggregateID != outboxEvent.AggregateID {
		t.Errorf("AggregateID mismatch: expected %s, got %s",
			outboxEvent.AggregateID, model.AggregateID)
	}

	// Verify default values
	if model.AttemptCount != 0 {
		t.Errorf("AttemptCount should be 0, got: %d", model.AttemptCount)
	}
	if model.LastError != "" {
		t.Errorf("LastError should be empty, got: %s", model.LastError)
	}
	if model.NextAttemptAt.IsZero() {
		t.Error("NextAttemptAt should not be zero")
	}
}

func TestEventModel_ToOutboxEvent(t *testing.T) {
	now := time.Now()
	model := &EventModel{
		ID:            uuid.New(),
		EventID:       uuid.New().String(),
		EventType:     "UserCreated",
		EventStatus:   string(abstraction.OutboxEventStatusPending),
		AggregateType: "User",
		AggregateID:   "user-123",
		Payload: map[string]interface{}{
			"email": "test@example.com",
		},
		AttemptCount:  1,
		LastError:     "some error",
		NextAttemptAt: now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	outboxEvent := model.ToOutboxEvent()

	// Verify all fields
	if outboxEvent.EventID != model.EventID {
		t.Errorf("EventID mismatch: expected %s, got %s", model.EventID, outboxEvent.EventID)
	}
	if outboxEvent.EventType != model.EventType {
		t.Errorf("EventType mismatch: expected %s, got %s", model.EventType, outboxEvent.EventType)
	}
	if outboxEvent.EventStatus != model.EventStatus {
		t.Errorf("EventStatus mismatch: expected %s, got %s", model.EventStatus, outboxEvent.EventStatus)
	}
	if outboxEvent.AggregateType != model.AggregateType {
		t.Errorf("AggregateType mismatch: expected %s, got %s",
			model.AggregateType, outboxEvent.AggregateType)
	}
	if outboxEvent.AggregateID != model.AggregateID {
		t.Errorf("AggregateID mismatch: expected %s, got %s",
			model.AggregateID, outboxEvent.AggregateID)
	}
}

func TestRegisterEvent(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}

func TestRegisterEvent_MultipleEvents(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}

func TestProcessBatch_SuccessfulEvents(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}

func TestProcessBatch_DifferentEventTypes(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}

func TestNewsqlEventChannel_DefaultValues(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}

func TestNewsqlEventChannel_CustomValues(t *testing.T) {
	t.Skip("Skipping database integration test - requires PostgreSQL")
}
