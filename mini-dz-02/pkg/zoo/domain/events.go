package domain

import (
	"time"
)

// Event is the base interface for all domain events
type Event interface {
	OccurredOn() time.Time
	EventType() string
}

// BaseEvent provides common functionality for all events
type BaseEvent struct {
	occurredOn time.Time
	eventType  string
}

// NewBaseEvent creates a new BaseEvent
func NewBaseEvent(eventType string) BaseEvent {
	return BaseEvent{
		occurredOn: time.Now(),
		eventType:  eventType,
	}
}

// OccurredOn returns when the event occurred
func (e BaseEvent) OccurredOn() time.Time {
	return e.occurredOn
}

// EventType returns the type of the event
func (e BaseEvent) EventType() string {
	return e.eventType
}

// AnimalMovedEvent represents an event when an animal is moved from one enclosure to another
type AnimalMovedEvent struct {
	BaseEvent
	AnimalID    string
	FromEnclosureID string
	ToEnclosureID   string
}

// NewAnimalMovedEvent creates a new AnimalMovedEvent
func NewAnimalMovedEvent(animalID, fromEnclosureID, toEnclosureID string) AnimalMovedEvent {
	return AnimalMovedEvent{
		BaseEvent:       NewBaseEvent("AnimalMoved"),
		AnimalID:        animalID,
		FromEnclosureID: fromEnclosureID,
		ToEnclosureID:   toEnclosureID,
	}
}

// FeedingTimeEvent represents an event when it's time to feed an animal
type FeedingTimeEvent struct {
	BaseEvent
	AnimalID string
	FoodType FoodType
}

// NewFeedingTimeEvent creates a new FeedingTimeEvent
func NewFeedingTimeEvent(animalID string, foodType FoodType) FeedingTimeEvent {
	return FeedingTimeEvent{
		BaseEvent: NewBaseEvent("FeedingTime"),
		AnimalID:  animalID,
		FoodType:  foodType,
	}
}