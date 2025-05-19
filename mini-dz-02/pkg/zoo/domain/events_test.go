package domain_test

import (
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

func TestBaseEvent(t *testing.T) {
	// Test creating a base event
	eventType := "TestEvent"
	baseEvent := domain.NewBaseEvent(eventType)

	// Check that the event type is set correctly
	if baseEvent.EventType() != eventType {
		t.Errorf("Expected event type to be %s, got %s", eventType, baseEvent.EventType())
	}

	// Check that the occurred on time is set to a reasonable value (within the last second)
	now := time.Now()
	occurredOn := baseEvent.OccurredOn()
	if occurredOn.After(now) || occurredOn.Before(now.Add(-1*time.Second)) {
		t.Errorf("Expected occurred on time to be within the last second, got %v", occurredOn)
	}
}

func TestAnimalMovedEvent(t *testing.T) {
	// Test creating an animal moved event
	animalID := "animal-1"
	fromEnclosureID := "enclosure-1"
	toEnclosureID := "enclosure-2"

	event := domain.NewAnimalMovedEvent(animalID, fromEnclosureID, toEnclosureID)

	// Check that the event type is set correctly
	if event.EventType() != "AnimalMoved" {
		t.Errorf("Expected event type to be AnimalMoved, got %s", event.EventType())
	}

	// Check that the animal ID is set correctly
	if event.AnimalID != animalID {
		t.Errorf("Expected animal ID to be %s, got %s", animalID, event.AnimalID)
	}

	// Check that the from enclosure ID is set correctly
	if event.FromEnclosureID != fromEnclosureID {
		t.Errorf("Expected from enclosure ID to be %s, got %s", fromEnclosureID, event.FromEnclosureID)
	}

	// Check that the to enclosure ID is set correctly
	if event.ToEnclosureID != toEnclosureID {
		t.Errorf("Expected to enclosure ID to be %s, got %s", toEnclosureID, event.ToEnclosureID)
	}

	// Check that the occurred on time is set to a reasonable value (within the last second)
	now := time.Now()
	occurredOn := event.OccurredOn()
	if occurredOn.After(now) || occurredOn.Before(now.Add(-1*time.Second)) {
		t.Errorf("Expected occurred on time to be within the last second, got %v", occurredOn)
	}
}

func TestFeedingTimeEvent(t *testing.T) {
	// Test creating a feeding time event
	animalID := "animal-1"
	foodType := domain.Meat

	event := domain.NewFeedingTimeEvent(animalID, foodType)

	// Check that the event type is set correctly
	if event.EventType() != "FeedingTime" {
		t.Errorf("Expected event type to be FeedingTime, got %s", event.EventType())
	}

	// Check that the animal ID is set correctly
	if event.AnimalID != animalID {
		t.Errorf("Expected animal ID to be %s, got %s", animalID, event.AnimalID)
	}

	// Check that the food type is set correctly
	if event.FoodType != foodType {
		t.Errorf("Expected food type to be %s, got %s", foodType, event.FoodType)
	}

	// Check that the occurred on time is set to a reasonable value (within the last second)
	now := time.Now()
	occurredOn := event.OccurredOn()
	if occurredOn.After(now) || occurredOn.Before(now.Add(-1*time.Second)) {
		t.Errorf("Expected occurred on time to be within the last second, got %v", occurredOn)
	}
}
