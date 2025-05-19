package domain_test

import (
	"fmt"
	"testing"

	"mini-dz-02/pkg/zoo/domain"
)

// Helper function to create a valid enclosure for testing
func createValidEnclosure(t *testing.T) *domain.Enclosure {
	size, err := domain.NewEnclosureSize(10)
	if err != nil {
		t.Fatalf("Failed to create enclosure size: %v", err)
	}

	capacity, err := domain.NewCapacity(5)
	if err != nil {
		t.Fatalf("Failed to create capacity: %v", err)
	}

	enclosure, err := domain.NewEnclosure(
		"enclosure-1",
		domain.Predator,
		size,
		capacity,
	)
	if err != nil {
		t.Fatalf("Failed to create enclosure: %v", err)
	}

	return enclosure
}

func TestNewEnclosure(t *testing.T) {
	// Test valid enclosure creation
	size, _ := domain.NewEnclosureSize(10)
	capacity, _ := domain.NewCapacity(5)

	enclosure, err := domain.NewEnclosure(
		"enclosure-1",
		domain.Predator,
		size,
		capacity,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if enclosure.ID != "enclosure-1" {
		t.Errorf("Expected ID to be 'enclosure-1', got %s", enclosure.ID)
	}

	if enclosure.Type != domain.Predator {
		t.Errorf("Expected Type to be Predator, got %s", enclosure.Type)
	}

	if enclosure.Size.Value() != 10 {
		t.Errorf("Expected Size to be 10, got %d", enclosure.Size.Value())
	}

	if enclosure.MaxCapacity.Value() != 5 {
		t.Errorf("Expected MaxCapacity to be 5, got %d", enclosure.MaxCapacity.Value())
	}

	if len(enclosure.CurrentAnimalIDs) != 0 {
		t.Errorf("Expected CurrentAnimalIDs to be empty, got %v", enclosure.CurrentAnimalIDs)
	}

	if !enclosure.CleaningStatus {
		t.Errorf("Expected CleaningStatus to be true (clean), got false")
	}

	// Test invalid enclosure creation (empty ID)
	_, err = domain.NewEnclosure(
		"",
		domain.Predator,
		size,
		capacity,
	)

	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}

	// Test invalid enclosure creation (size < capacity)
	smallSize, _ := domain.NewEnclosureSize(3)
	largeCapacity, _ := domain.NewCapacity(5)

	_, err = domain.NewEnclosure(
		"enclosure-2",
		domain.Predator,
		smallSize,
		largeCapacity,
	)

	if err == nil {
		t.Error("Expected error for size < capacity, got nil")
	}
}

func TestEnclosure_AddAnimal(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// Test adding a valid animal
	err := enclosure.AddAnimal("animal-1")
	if err != nil {
		t.Errorf("Expected no error when adding a valid animal, got %v", err)
	}

	if len(enclosure.CurrentAnimalIDs) != 1 {
		t.Errorf("Expected 1 animal in enclosure, got %d", len(enclosure.CurrentAnimalIDs))
	}

	if enclosure.CurrentAnimalIDs[0] != "animal-1" {
		t.Errorf("Expected animal ID to be 'animal-1', got %s", enclosure.CurrentAnimalIDs[0])
	}

	if enclosure.CleaningStatus {
		t.Errorf("Expected CleaningStatus to be false (dirty) after adding an animal, got true")
	}

	// Test adding an empty animal ID
	err = enclosure.AddAnimal("")
	if err == nil {
		t.Error("Expected error when adding an empty animal ID, got nil")
	}

	// Test adding the same animal again
	err = enclosure.AddAnimal("animal-1")
	if err == nil {
		t.Error("Expected error when adding the same animal again, got nil")
	}

	// Test adding animals until capacity is reached
	enclosure.AddAnimal("animal-2")
	enclosure.AddAnimal("animal-3")
	enclosure.AddAnimal("animal-4")
	enclosure.AddAnimal("animal-5")

	// Now the enclosure should be at capacity (5 animals)
	err = enclosure.AddAnimal("animal-6")
	if err == nil {
		t.Error("Expected error when adding an animal to a full enclosure, got nil")
	}
}

func TestEnclosure_RemoveAnimal(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// Add some animals
	enclosure.AddAnimal("animal-1")
	enclosure.AddAnimal("animal-2")
	enclosure.AddAnimal("animal-3")

	// Test removing a valid animal
	err := enclosure.RemoveAnimal("animal-2")
	if err != nil {
		t.Errorf("Expected no error when removing a valid animal, got %v", err)
	}

	if len(enclosure.CurrentAnimalIDs) != 2 {
		t.Errorf("Expected 2 animals in enclosure after removal, got %d", len(enclosure.CurrentAnimalIDs))
	}

	// Check that the correct animal was removed
	for _, id := range enclosure.CurrentAnimalIDs {
		if id == "animal-2" {
			t.Errorf("Expected animal-2 to be removed, but it's still in the enclosure")
		}
	}

	// Test removing an empty animal ID
	err = enclosure.RemoveAnimal("")
	if err == nil {
		t.Error("Expected error when removing an empty animal ID, got nil")
	}

	// Test removing an animal that's not in the enclosure
	err = enclosure.RemoveAnimal("animal-4")
	if err == nil {
		t.Error("Expected error when removing an animal that's not in the enclosure, got nil")
	}
}

func TestEnclosure_Clean(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// New enclosures are clean, so cleaning should fail
	err := enclosure.Clean()
	if err == nil {
		t.Error("Expected error when cleaning an already clean enclosure, got nil")
	}

	// Make the enclosure dirty by adding an animal
	enclosure.AddAnimal("animal-1")

	// Now cleaning should succeed
	err = enclosure.Clean()
	if err != nil {
		t.Errorf("Expected no error when cleaning a dirty enclosure, got %v", err)
	}

	if !enclosure.CleaningStatus {
		t.Errorf("Expected CleaningStatus to be true (clean) after cleaning, got false")
	}

	// Cleaning again should fail
	err = enclosure.Clean()
	if err == nil {
		t.Error("Expected error when cleaning an already clean enclosure, got nil")
	}
}

func TestEnclosure_IsClean(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// New enclosures are clean
	if !enclosure.IsClean() {
		t.Error("Expected new enclosure to be clean")
	}

	// Make the enclosure dirty by adding an animal
	enclosure.AddAnimal("animal-1")

	if enclosure.IsClean() {
		t.Error("Expected enclosure to be dirty after adding an animal")
	}

	// Clean the enclosure
	enclosure.Clean()

	if !enclosure.IsClean() {
		t.Error("Expected enclosure to be clean after cleaning")
	}
}

func TestEnclosure_CurrentAnimalCount(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// New enclosures have no animals
	if enclosure.CurrentAnimalCount() != 0 {
		t.Errorf("Expected new enclosure to have 0 animals, got %d", enclosure.CurrentAnimalCount())
	}

	// Add some animals
	enclosure.AddAnimal("animal-1")
	enclosure.AddAnimal("animal-2")

	if enclosure.CurrentAnimalCount() != 2 {
		t.Errorf("Expected enclosure to have 2 animals after adding, got %d", enclosure.CurrentAnimalCount())
	}

	// Remove an animal
	enclosure.RemoveAnimal("animal-1")

	if enclosure.CurrentAnimalCount() != 1 {
		t.Errorf("Expected enclosure to have 1 animal after removing, got %d", enclosure.CurrentAnimalCount())
	}
}

func TestEnclosure_HasSpace(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// New enclosures have space
	if !enclosure.HasSpace() {
		t.Error("Expected new enclosure to have space")
	}

	// Add animals until capacity is reached
	for i := 1; i <= 5; i++ {
		enclosure.AddAnimal(fmt.Sprintf("animal-%d", i))
	}

	// Now the enclosure should be full
	if enclosure.HasSpace() {
		t.Error("Expected enclosure to be full after adding 5 animals")
	}

	// Remove an animal
	enclosure.RemoveAnimal("animal-1")

	// Now the enclosure should have space again
	if !enclosure.HasSpace() {
		t.Error("Expected enclosure to have space after removing an animal")
	}
}

func TestEnclosure_IsEmpty(t *testing.T) {
	enclosure := createValidEnclosure(t)

	// New enclosures are empty
	if !enclosure.IsEmpty() {
		t.Error("Expected new enclosure to be empty")
	}

	// Add an animal
	enclosure.AddAnimal("animal-1")

	// Now the enclosure should not be empty
	if enclosure.IsEmpty() {
		t.Error("Expected enclosure to not be empty after adding an animal")
	}

	// Remove the animal
	enclosure.RemoveAnimal("animal-1")

	// Now the enclosure should be empty again
	if !enclosure.IsEmpty() {
		t.Error("Expected enclosure to be empty after removing all animals")
	}
}
