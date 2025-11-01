//go:build integration
// +build integration

package sqlchannel

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ali63yavari/outbox_abstraction/abstraction"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Integration tests require a real PostgreSQL database
// Run with: go test -tags=integration -v ./sql_channel
//
// Set the following environment variables:
// export TEST_DB_HOST=localhost
// export TEST_DB_PORT=5432
// export TEST_DB_USER=postgres
// export TEST_DB_PASSWORD=postgres
// export TEST_DB_NAME=outbox_test

func setupIntegrationDB(t *testing.T) *gorm.DB {
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	dbname := os.Getenv("TEST_DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}
	if dbname == "" {
		dbname = "outbox_test"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Skipf("Failed to connect to PostgreSQL: %v. Set TEST_DB_* env vars or skip integration tests", err)
	}

	// Clean up and migrate
	db.Exec("DROP TABLE IF EXISTS event_models CASCADE")
	err = db.AutoMigrate(&EventModel{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestIntegration_RegisterEvent(t *testing.T) {
	db := setupIntegrationDB(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eventType := mockEventType{name: "UserCreated", id: 1}
	handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		return nil
	}

	interval := 10 * time.Second
	maxRetries := 3
	batchSize := 10

	channel := NewsqlEventChannel(ctx, db, eventType, &maxRetries, &batchSize, &interval, handler)

	// Create an event
	event := abstraction.CreateNewEvent(eventType, "User", "user-123", map[string]interface{}{
		"email": "test@example.com",
	})

	// Register the event
	err := channel.RegisterEvent(event)
	if err != nil {
		t.Fatalf("Failed to register event: %v", err)
	}

	// Verify event was saved to database
	var count int64
	db.Model(&EventModel{}).Where("event_id = ?", event.EventID).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 event in database, got: %d", count)
	}

	// Verify event details
	var savedEvent EventModel
	err = db.Where("event_id = ?", event.EventID).First(&savedEvent).Error
	if err != nil {
		t.Fatalf("Failed to retrieve saved event: %v", err)
	}

	if savedEvent.EventType != "UserCreated" {
		t.Errorf("Expected EventType 'UserCreated', got: %s", savedEvent.EventType)
	}
	if savedEvent.EventStatus != string(abstraction.OutboxEventStatusPending) {
		t.Errorf("Expected EventStatus 'pending', got: %s", savedEvent.EventStatus)
	}
}

func TestIntegration_ProcessBatch_SuccessfulEvents(t *testing.T) {
	db := setupIntegrationDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	eventType := mockEventType{name: "UserCreated", id: 1}
	processedEvents := []string{}

	handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		processedEvents = append(processedEvents, event.EventID)
		return nil
	}

	interval := 200 * time.Millisecond
	maxRetries := 3
	batchSize := 10

	channel := NewsqlEventChannel(ctx, db, eventType, &maxRetries, &batchSize, &interval, handler)

	// Create and register events
	event1 := abstraction.CreateNewEvent(eventType, "User", "user-1", map[string]interface{}{})
	event2 := abstraction.CreateNewEvent(eventType, "User", "user-2", map[string]interface{}{})

	channel.RegisterEvent(event1)
	channel.RegisterEvent(event2)

	// Wait for processing
	time.Sleep(600 * time.Millisecond)
	cancel()

	// Verify events were processed
	if len(processedEvents) < 2 {
		t.Errorf("Expected at least 2 events to be processed, got: %d", len(processedEvents))
	}

	// Verify events are marked as closed
	var closedCount int64
	db.Model(&EventModel{}).
		Where("event_status = ?", string(abstraction.OutboxEventStatusClosed)).
		Count(&closedCount)

	if closedCount < 2 {
		t.Errorf("Expected at least 2 closed events, got: %d", closedCount)
	}
}

func TestIntegration_ProcessBatch_DifferentEventTypes(t *testing.T) {
	db := setupIntegrationDB(t)
	ctx1, cancel1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel1()
	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()

	eventType1 := mockEventType{name: "UserCreated", id: 1}
	eventType2 := mockEventType{name: "UserUpdated", id: 2}

	processedType1 := []string{}
	processedType2 := []string{}

	handler1 := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		processedType1 = append(processedType1, event.EventID)
		return nil
	}

	handler2 := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		processedType2 = append(processedType2, event.EventID)
		return nil
	}

	interval := 200 * time.Millisecond
	maxRetries := 3
	batchSize := 10

	channel1 := NewsqlEventChannel(ctx1, db, eventType1, &maxRetries, &batchSize, &interval, handler1)
	channel2 := NewsqlEventChannel(ctx2, db, eventType2, &maxRetries, &batchSize, &interval, handler2)

	// Create events of different types
	event1 := abstraction.CreateNewEvent(eventType1, "User", "user-1", map[string]interface{}{})
	event2 := abstraction.CreateNewEvent(eventType2, "User", "user-2", map[string]interface{}{})

	channel1.RegisterEvent(event1)
	channel2.RegisterEvent(event2)

	// Wait for processing
	time.Sleep(600 * time.Millisecond)
	cancel1()
	cancel2()

	// Verify each channel processed only its events
	if len(processedType1) < 1 {
		t.Errorf("Expected channel 1 to process at least 1 event, got: %d", len(processedType1))
	}
	if len(processedType2) < 1 {
		t.Errorf("Expected channel 2 to process at least 1 event, got: %d", len(processedType2))
	}

	if len(processedType1) > 0 && processedType1[0] != event1.EventID {
		t.Error("Channel 1 should process event1")
	}
	if len(processedType2) > 0 && processedType2[0] != event2.EventID {
		t.Error("Channel 2 should process event2")
	}
}

func TestIntegration_FailedEvent_Retry(t *testing.T) {
	db := setupIntegrationDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventType := mockEventType{name: "FailingEvent", id: 1}
	attempts := 0

	handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		attempts++
		if attempts < 3 {
			return fmt.Errorf("simulated failure %d", attempts)
		}
		return nil
	}

	interval := 200 * time.Millisecond
	maxRetries := 5
	batchSize := 10

	channel := NewsqlEventChannel(ctx, db, eventType, &maxRetries, &batchSize, &interval, handler)

	// Create and register event
	event := abstraction.CreateNewEvent(eventType, "Test", "test-1", map[string]interface{}{})
	channel.RegisterEvent(event)

	// Wait for retries
	time.Sleep(2 * time.Second)
	cancel()

	// Verify event was eventually processed
	if attempts < 3 {
		t.Errorf("Expected at least 3 attempts, got: %d", attempts)
	}

	// Verify event is marked as closed
	var savedEvent EventModel
	err := db.Where("event_id = ?", event.EventID).First(&savedEvent).Error
	if err != nil {
		t.Fatalf("Failed to retrieve event: %v", err)
	}

	if savedEvent.EventStatus != string(abstraction.OutboxEventStatusClosed) {
		t.Errorf("Expected event status 'closed', got: %s", savedEvent.EventStatus)
	}
}

func TestIntegration_MaxRetries_MarkAsFailed(t *testing.T) {
	db := setupIntegrationDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventType := mockEventType{name: "AlwaysFailingEvent", id: 1}

	handler := func(ctx context.Context, event *abstraction.OutboxEvent) error {
		return fmt.Errorf("always fails")
	}

	interval := 100 * time.Millisecond
	maxRetries := 2
	batchSize := 10

	channel := NewsqlEventChannel(ctx, db, eventType, &maxRetries, &batchSize, &interval, handler)

	// Create and register event
	event := abstraction.CreateNewEvent(eventType, "Test", "test-1", map[string]interface{}{})
	channel.RegisterEvent(event)

	// Wait for retries
	time.Sleep(2 * time.Second)
	cancel()

	// Verify event is marked as failed
	var savedEvent EventModel
	err := db.Where("event_id = ?", event.EventID).First(&savedEvent).Error
	if err != nil {
		t.Fatalf("Failed to retrieve event: %v", err)
	}

	if savedEvent.EventStatus != string(abstraction.OutboxEventStatusFailed) {
		t.Errorf("Expected event status 'failed', got: %s", savedEvent.EventStatus)
	}

	if savedEvent.AttemptCount < maxRetries {
		t.Errorf("Expected at least %d attempts, got: %d", maxRetries, savedEvent.AttemptCount)
	}
}
