package infrastructure_test

import (
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/domain"
	"mini-dz-02/pkg/zoo/infrastructure"
)

// Helper function to create a valid animal for testing
func createTestAnimal(id string) *domain.Animal {
	species, _ := domain.NewSpecies("Lion")
	name, _ := domain.NewAnimalName("Leo")
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-5, 0, 0))

	animal, _ := domain.NewAnimal(
		id,
		species,
		name,
		birthDate,
		domain.Male,
		domain.Meat,
		domain.Healthy,
	)

	return animal
}

// Helper function to create a valid enclosure for testing
func createTestEnclosure(id string) *domain.Enclosure {
	size, _ := domain.NewEnclosureSize(10)
	capacity, _ := domain.NewCapacity(5)

	enclosure, _ := domain.NewEnclosure(
		id,
		domain.Predator,
		size,
		capacity,
	)

	return enclosure
}

// Helper function to create a valid feeding schedule for testing
func createTestFeedingSchedule(id string, animalID string) *domain.FeedingSchedule {
	feedingTime, _ := domain.NewFeedingTime(time.Now().Add(1 * time.Hour))

	schedule, _ := domain.NewFeedingSchedule(
		id,
		animalID,
		feedingTime,
		domain.Meat,
	)

	return schedule
}

func TestInMemoryAnimalRepository(t *testing.T) {
	repo := infrastructure.NewInMemoryAnimalRepository()

	// Test GetByID with non-existent animal
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent animal, got nil")
	}

	// Test GetAll with empty repository
	animals, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all animals from empty repository, got %v", err)
	}
	if len(animals) != 0 {
		t.Errorf("Expected empty slice when getting all animals from empty repository, got %d animals", len(animals))
	}

	// Test Save with nil animal
	err = repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil animal, got nil")
	}

	// Test Save and GetByID
	animal := createTestAnimal("animal-1")
	err = repo.Save(animal)
	if err != nil {
		t.Errorf("Expected no error when saving animal, got %v", err)
	}

	retrievedAnimal, err := repo.GetByID("animal-1")
	if err != nil {
		t.Errorf("Expected no error when getting animal, got %v", err)
	}
	if retrievedAnimal.ID != "animal-1" {
		t.Errorf("Expected animal ID to be 'animal-1', got %s", retrievedAnimal.ID)
	}

	// Test GetAll with non-empty repository
	animals, err = repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all animals, got %v", err)
	}
	if len(animals) != 1 {
		t.Errorf("Expected 1 animal, got %d", len(animals))
	}

	// Test Delete with non-existent animal
	err = repo.Delete("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent animal, got nil")
	}

	// Test Delete with existing animal
	err = repo.Delete("animal-1")
	if err != nil {
		t.Errorf("Expected no error when deleting animal, got %v", err)
	}

	// Verify animal was deleted
	_, err = repo.GetByID("animal-1")
	if err == nil {
		t.Error("Expected error when getting deleted animal, got nil")
	}

	// Test GetByEnclosureID
	animal1 := createTestAnimal("animal-1")
	animal1.EnclosureID = "enclosure-1"
	animal2 := createTestAnimal("animal-2")
	animal2.EnclosureID = "enclosure-1"
	animal3 := createTestAnimal("animal-3")
	animal3.EnclosureID = "enclosure-2"

	repo.Save(animal1)
	repo.Save(animal2)
	repo.Save(animal3)

	animalsInEnclosure1, err := repo.GetByEnclosureID("enclosure-1")
	if err != nil {
		t.Errorf("Expected no error when getting animals by enclosure ID, got %v", err)
	}
	if len(animalsInEnclosure1) != 2 {
		t.Errorf("Expected 2 animals in enclosure-1, got %d", len(animalsInEnclosure1))
	}

	animalsInEnclosure2, err := repo.GetByEnclosureID("enclosure-2")
	if err != nil {
		t.Errorf("Expected no error when getting animals by enclosure ID, got %v", err)
	}
	if len(animalsInEnclosure2) != 1 {
		t.Errorf("Expected 1 animal in enclosure-2, got %d", len(animalsInEnclosure2))
	}

	animalsInEnclosure3, err := repo.GetByEnclosureID("enclosure-3")
	if err != nil {
		t.Errorf("Expected no error when getting animals by non-existent enclosure ID, got %v", err)
	}
	if len(animalsInEnclosure3) != 0 {
		t.Errorf("Expected 0 animals in non-existent enclosure, got %d", len(animalsInEnclosure3))
	}
}

func TestInMemoryEnclosureRepository(t *testing.T) {
	repo := infrastructure.NewInMemoryEnclosureRepository()

	// Test GetByID with non-existent enclosure
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent enclosure, got nil")
	}

	// Test GetAll with empty repository
	enclosures, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all enclosures from empty repository, got %v", err)
	}
	if len(enclosures) != 0 {
		t.Errorf("Expected empty slice when getting all enclosures from empty repository, got %d enclosures", len(enclosures))
	}

	// Test Save with nil enclosure
	err = repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil enclosure, got nil")
	}

	// Test Save and GetByID
	enclosure := createTestEnclosure("enclosure-1")
	err = repo.Save(enclosure)
	if err != nil {
		t.Errorf("Expected no error when saving enclosure, got %v", err)
	}

	retrievedEnclosure, err := repo.GetByID("enclosure-1")
	if err != nil {
		t.Errorf("Expected no error when getting enclosure, got %v", err)
	}
	if retrievedEnclosure.ID != "enclosure-1" {
		t.Errorf("Expected enclosure ID to be 'enclosure-1', got %s", retrievedEnclosure.ID)
	}

	// Test GetAll with non-empty repository
	enclosures, err = repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all enclosures, got %v", err)
	}
	if len(enclosures) != 1 {
		t.Errorf("Expected 1 enclosure, got %d", len(enclosures))
	}

	// Test Delete with non-existent enclosure
	err = repo.Delete("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent enclosure, got nil")
	}

	// Test Delete with existing enclosure
	err = repo.Delete("enclosure-1")
	if err != nil {
		t.Errorf("Expected no error when deleting enclosure, got %v", err)
	}

	// Verify enclosure was deleted
	_, err = repo.GetByID("enclosure-1")
	if err == nil {
		t.Error("Expected error when getting deleted enclosure, got nil")
	}

	// Test GetByType
	enclosure1 := createTestEnclosure("enclosure-1")
	enclosure1.Type = domain.Predator
	enclosure2 := createTestEnclosure("enclosure-2")
	enclosure2.Type = domain.Predator
	enclosure3 := createTestEnclosure("enclosure-3")
	enclosure3.Type = domain.Herbivore

	repo.Save(enclosure1)
	repo.Save(enclosure2)
	repo.Save(enclosure3)

	predatorEnclosures, err := repo.GetByType(domain.Predator)
	if err != nil {
		t.Errorf("Expected no error when getting enclosures by type, got %v", err)
	}
	if len(predatorEnclosures) != 2 {
		t.Errorf("Expected 2 predator enclosures, got %d", len(predatorEnclosures))
	}

	herbivoreEnclosures, err := repo.GetByType(domain.Herbivore)
	if err != nil {
		t.Errorf("Expected no error when getting enclosures by type, got %v", err)
	}
	if len(herbivoreEnclosures) != 1 {
		t.Errorf("Expected 1 herbivore enclosure, got %d", len(herbivoreEnclosures))
	}

	aviaryEnclosures, err := repo.GetByType(domain.Aviary)
	if err != nil {
		t.Errorf("Expected no error when getting enclosures by non-existent type, got %v", err)
	}
	if len(aviaryEnclosures) != 0 {
		t.Errorf("Expected 0 aviary enclosures, got %d", len(aviaryEnclosures))
	}

	// Test GetAvailable
	// Make enclosure1 full
	for i := 0; i < 5; i++ {
		enclosure1.AddAnimal("animal-" + string(rune('a'+i)))
	}
	repo.Save(enclosure1)

	availableEnclosures, err := repo.GetAvailable()
	if err != nil {
		t.Errorf("Expected no error when getting available enclosures, got %v", err)
	}
	if len(availableEnclosures) != 2 {
		t.Errorf("Expected 2 available enclosures, got %d", len(availableEnclosures))
	}
}

func TestInMemoryFeedingScheduleRepository(t *testing.T) {
	repo := infrastructure.NewInMemoryFeedingScheduleRepository()

	// Test GetByID with non-existent schedule
	_, err := repo.GetByID("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent schedule, got nil")
	}

	// Test GetAll with empty repository
	schedules, err := repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all schedules from empty repository, got %v", err)
	}
	if len(schedules) != 0 {
		t.Errorf("Expected empty slice when getting all schedules from empty repository, got %d schedules", len(schedules))
	}

	// Test Save with nil schedule
	err = repo.Save(nil)
	if err == nil {
		t.Error("Expected error when saving nil schedule, got nil")
	}

	// Test Save and GetByID
	schedule := createTestFeedingSchedule("schedule-1", "animal-1")
	err = repo.Save(schedule)
	if err != nil {
		t.Errorf("Expected no error when saving schedule, got %v", err)
	}

	retrievedSchedule, err := repo.GetByID("schedule-1")
	if err != nil {
		t.Errorf("Expected no error when getting schedule, got %v", err)
	}
	if retrievedSchedule.ID != "schedule-1" {
		t.Errorf("Expected schedule ID to be 'schedule-1', got %s", retrievedSchedule.ID)
	}

	// Test GetAll with non-empty repository
	schedules, err = repo.GetAll()
	if err != nil {
		t.Errorf("Expected no error when getting all schedules, got %v", err)
	}
	if len(schedules) != 1 {
		t.Errorf("Expected 1 schedule, got %d", len(schedules))
	}

	// Test Delete with non-existent schedule
	err = repo.Delete("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent schedule, got nil")
	}

	// Test Delete with existing schedule
	err = repo.Delete("schedule-1")
	if err != nil {
		t.Errorf("Expected no error when deleting schedule, got %v", err)
	}

	// Verify schedule was deleted
	_, err = repo.GetByID("schedule-1")
	if err == nil {
		t.Error("Expected error when getting deleted schedule, got nil")
	}

	// Test GetByAnimalID
	schedule1 := createTestFeedingSchedule("schedule-1", "animal-1")
	schedule2 := createTestFeedingSchedule("schedule-2", "animal-1")
	schedule3 := createTestFeedingSchedule("schedule-3", "animal-2")

	repo.Save(schedule1)
	repo.Save(schedule2)
	repo.Save(schedule3)

	schedulesForAnimal1, err := repo.GetByAnimalID("animal-1")
	if err != nil {
		t.Errorf("Expected no error when getting schedules by animal ID, got %v", err)
	}
	if len(schedulesForAnimal1) != 2 {
		t.Errorf("Expected 2 schedules for animal-1, got %d", len(schedulesForAnimal1))
	}

	schedulesForAnimal2, err := repo.GetByAnimalID("animal-2")
	if err != nil {
		t.Errorf("Expected no error when getting schedules by animal ID, got %v", err)
	}
	if len(schedulesForAnimal2) != 1 {
		t.Errorf("Expected 1 schedule for animal-2, got %d", len(schedulesForAnimal2))
	}

	schedulesForAnimal3, err := repo.GetByAnimalID("animal-3")
	if err != nil {
		t.Errorf("Expected no error when getting schedules by non-existent animal ID, got %v", err)
	}
	if len(schedulesForAnimal3) != 0 {
		t.Errorf("Expected 0 schedules for non-existent animal, got %d", len(schedulesForAnimal3))
	}

	// Test GetDueSchedules
	// Create a past schedule
	pastFeedingTime, _ := domain.NewFeedingTime(time.Now().Add(-1 * time.Hour))
	pastSchedule, _ := domain.NewFeedingSchedule("schedule-past", "animal-1", pastFeedingTime, domain.Meat)
	repo.Save(pastSchedule)

	// Create a completed past schedule
	completedPastSchedule, _ := domain.NewFeedingSchedule("schedule-completed", "animal-1", pastFeedingTime, domain.Meat)
	completedPastSchedule.MarkCompleted()
	repo.Save(completedPastSchedule)

	dueSchedules, err := repo.GetDueSchedules()
	if err != nil {
		t.Errorf("Expected no error when getting due schedules, got %v", err)
	}
	if len(dueSchedules) != 1 {
		t.Errorf("Expected 1 due schedule, got %d", len(dueSchedules))
	}
	if dueSchedules[0].ID != "schedule-past" {
		t.Errorf("Expected due schedule ID to be 'schedule-past', got %s", dueSchedules[0].ID)
	}
}

func TestInMemoryEventPublisher(t *testing.T) {
	publisher := infrastructure.NewInMemoryEventPublisher()

	// Test Publish with nil event
	err := publisher.Publish(nil)
	if err == nil {
		t.Error("Expected error when publishing nil event, got nil")
	}

	// Test GetEvents with empty publisher
	events := publisher.GetEvents()
	if len(events) != 0 {
		t.Errorf("Expected empty slice when getting events from empty publisher, got %d events", len(events))
	}

	// Test Publish and GetEvents
	animalMovedEvent := domain.NewAnimalMovedEvent("animal-1", "enclosure-1", "enclosure-2")
	err = publisher.Publish(animalMovedEvent)
	if err != nil {
		t.Errorf("Expected no error when publishing event, got %v", err)
	}

	events = publisher.GetEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	if events[0].EventType() != "AnimalMoved" {
		t.Errorf("Expected event type to be 'AnimalMoved', got %s", events[0].EventType())
	}

	// Test publishing multiple events
	feedingTimeEvent := domain.NewFeedingTimeEvent("animal-1", domain.Meat)
	err = publisher.Publish(feedingTimeEvent)
	if err != nil {
		t.Errorf("Expected no error when publishing event, got %v", err)
	}

	events = publisher.GetEvents()
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}
