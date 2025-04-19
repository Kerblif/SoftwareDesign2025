package domain_test

import (
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

// Helper function to create a valid animal for testing
func createValidAnimal(t *testing.T) *domain.Animal {
	species, err := domain.NewSpecies("Lion")
	if err != nil {
		t.Fatalf("Failed to create species: %v", err)
	}

	name, err := domain.NewAnimalName("Leo")
	if err != nil {
		t.Fatalf("Failed to create animal name: %v", err)
	}

	birthDate, err := domain.NewBirthDate(time.Now().AddDate(-5, 0, 0)) // 5 years ago
	if err != nil {
		t.Fatalf("Failed to create birth date: %v", err)
	}

	animal, err := domain.NewAnimal(
		"animal-1",
		species,
		name,
		birthDate,
		domain.Male,
		domain.Meat,
		domain.Healthy,
	)
	if err != nil {
		t.Fatalf("Failed to create animal: %v", err)
	}

	return animal
}

func TestNewAnimal(t *testing.T) {
	// Test valid animal creation
	species, _ := domain.NewSpecies("Lion")
	name, _ := domain.NewAnimalName("Leo")
	birthDate, _ := domain.NewBirthDate(time.Now().AddDate(-5, 0, 0))

	animal, err := domain.NewAnimal(
		"animal-1",
		species,
		name,
		birthDate,
		domain.Male,
		domain.Meat,
		domain.Healthy,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if animal.ID != "animal-1" {
		t.Errorf("Expected ID to be 'animal-1', got %s", animal.ID)
	}

	if animal.Species.String() != "Lion" {
		t.Errorf("Expected Species to be 'Lion', got %s", animal.Species.String())
	}

	if animal.Name.String() != "Leo" {
		t.Errorf("Expected Name to be 'Leo', got %s", animal.Name.String())
	}

	if animal.Gender != domain.Male {
		t.Errorf("Expected Gender to be Male, got %s", animal.Gender)
	}

	if animal.FavoriteFood != domain.Meat {
		t.Errorf("Expected FavoriteFood to be Meat, got %s", animal.FavoriteFood)
	}

	if animal.HealthStatus != domain.Healthy {
		t.Errorf("Expected HealthStatus to be Healthy, got %s", animal.HealthStatus)
	}

	// Test invalid animal creation (empty ID)
	_, err = domain.NewAnimal(
		"",
		species,
		name,
		birthDate,
		domain.Male,
		domain.Meat,
		domain.Healthy,
	)

	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}
}

func TestAnimal_Feed(t *testing.T) {
	// Test feeding a healthy animal
	animal := createValidAnimal(t)
	err := animal.Feed(domain.Meat)
	if err != nil {
		t.Errorf("Expected no error when feeding a healthy animal, got %v", err)
	}

	// Test feeding a sick animal
	sickAnimal := createValidAnimal(t)
	sickAnimal.SetHealthStatus(domain.Sick)
	err = sickAnimal.Feed(domain.Meat)
	if err == nil {
		t.Error("Expected error when feeding a sick animal, got nil")
	}
}

func TestAnimal_Treat(t *testing.T) {
	// Test treating a sick animal
	animal := createValidAnimal(t)
	animal.SetHealthStatus(domain.Sick)
	err := animal.Treat()
	if err != nil {
		t.Errorf("Expected no error when treating a sick animal, got %v", err)
	}
	if animal.HealthStatus != domain.Healthy {
		t.Errorf("Expected animal to be healthy after treatment, got %s", animal.HealthStatus)
	}

	// Test treating a healthy animal
	healthyAnimal := createValidAnimal(t)
	err = healthyAnimal.Treat()
	if err == nil {
		t.Error("Expected error when treating a healthy animal, got nil")
	}
}

func TestAnimal_MoveToEnclosure(t *testing.T) {
	// Test moving an animal to a valid enclosure
	animal := createValidAnimal(t)
	err := animal.MoveToEnclosure("enclosure-1")
	if err != nil {
		t.Errorf("Expected no error when moving to a valid enclosure, got %v", err)
	}
	if animal.EnclosureID != "enclosure-1" {
		t.Errorf("Expected EnclosureID to be 'enclosure-1', got %s", animal.EnclosureID)
	}

	// Test moving an animal to an empty enclosure ID
	err = animal.MoveToEnclosure("")
	if err == nil {
		t.Error("Expected error when moving to an empty enclosure ID, got nil")
	}

	// Test moving an animal to the same enclosure
	err = animal.MoveToEnclosure("enclosure-1")
	if err == nil {
		t.Error("Expected error when moving to the same enclosure, got nil")
	}
}

func TestAnimal_SetHealthStatus(t *testing.T) {
	animal := createValidAnimal(t)

	// Test setting to sick
	animal.SetHealthStatus(domain.Sick)
	if animal.HealthStatus != domain.Sick {
		t.Errorf("Expected HealthStatus to be Sick, got %s", animal.HealthStatus)
	}

	// Test setting back to healthy
	animal.SetHealthStatus(domain.Healthy)
	if animal.HealthStatus != domain.Healthy {
		t.Errorf("Expected HealthStatus to be Healthy, got %s", animal.HealthStatus)
	}
}

func TestAnimal_IsInEnclosure(t *testing.T) {
	// Test animal not in enclosure
	animal := createValidAnimal(t)
	if animal.IsInEnclosure() {
		t.Error("Expected animal to not be in an enclosure")
	}

	// Test animal in enclosure
	animal.MoveToEnclosure("enclosure-1")
	if !animal.IsInEnclosure() {
		t.Error("Expected animal to be in an enclosure")
	}
}
