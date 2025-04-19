package domain_test

import (
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

func TestSpecies(t *testing.T) {
	// Test valid species creation
	validSpecies := "Lion"
	species, err := domain.NewSpecies(validSpecies)
	if err != nil {
		t.Errorf("Expected no error for valid species, got %v", err)
	}
	if species.String() != validSpecies {
		t.Errorf("Expected species string to be %s, got %s", validSpecies, species.String())
	}

	// Test invalid species creation (empty)
	_, err = domain.NewSpecies("")
	if err == nil {
		t.Error("Expected error for empty species, got nil")
	}
}

func TestAnimalName(t *testing.T) {
	// Test valid animal name creation
	validName := "Leo"
	name, err := domain.NewAnimalName(validName)
	if err != nil {
		t.Errorf("Expected no error for valid name, got %v", err)
	}
	if name.String() != validName {
		t.Errorf("Expected name string to be %s, got %s", validName, name.String())
	}

	// Test invalid animal name creation (empty)
	_, err = domain.NewAnimalName("")
	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestBirthDate(t *testing.T) {
	// Test valid birth date creation (in the past)
	validDate := time.Now().AddDate(-5, 0, 0) // 5 years ago
	birthDate, err := domain.NewBirthDate(validDate)
	if err != nil {
		t.Errorf("Expected no error for valid birth date, got %v", err)
	}
	if !birthDate.Time().Equal(validDate) {
		t.Errorf("Expected birth date to be %v, got %v", validDate, birthDate.Time())
	}

	// Test invalid birth date creation (in the future)
	futureDate := time.Now().AddDate(1, 0, 0) // 1 year in the future
	_, err = domain.NewBirthDate(futureDate)
	if err == nil {
		t.Error("Expected error for future birth date, got nil")
	}
}

func TestEnclosureSize(t *testing.T) {
	// Test valid enclosure size creation
	validSize := 10
	size, err := domain.NewEnclosureSize(validSize)
	if err != nil {
		t.Errorf("Expected no error for valid size, got %v", err)
	}
	if size.Value() != validSize {
		t.Errorf("Expected size value to be %d, got %d", validSize, size.Value())
	}

	// Test invalid enclosure size creation (zero)
	_, err = domain.NewEnclosureSize(0)
	if err == nil {
		t.Error("Expected error for zero size, got nil")
	}

	// Test invalid enclosure size creation (negative)
	_, err = domain.NewEnclosureSize(-1)
	if err == nil {
		t.Error("Expected error for negative size, got nil")
	}
}

func TestCapacity(t *testing.T) {
	// Test valid capacity creation
	validCapacity := 5
	capacity, err := domain.NewCapacity(validCapacity)
	if err != nil {
		t.Errorf("Expected no error for valid capacity, got %v", err)
	}
	if capacity.Value() != validCapacity {
		t.Errorf("Expected capacity value to be %d, got %d", validCapacity, capacity.Value())
	}

	// Test invalid capacity creation (zero)
	_, err = domain.NewCapacity(0)
	if err == nil {
		t.Error("Expected error for zero capacity, got nil")
	}

	// Test invalid capacity creation (negative)
	_, err = domain.NewCapacity(-1)
	if err == nil {
		t.Error("Expected error for negative capacity, got nil")
	}
}

func TestFeedingTime(t *testing.T) {
	// Test valid feeding time creation
	validTime := time.Now().Add(1 * time.Hour)
	feedingTime, err := domain.NewFeedingTime(validTime)
	if err != nil {
		t.Errorf("Expected no error for valid feeding time, got %v", err)
	}
	if !feedingTime.Time().Equal(validTime) {
		t.Errorf("Expected feeding time to be %v, got %v", validTime, feedingTime.Time())
	}

	// Test feeding time in the past (should still be valid)
	pastTime := time.Now().Add(-1 * time.Hour)
	feedingTime, err = domain.NewFeedingTime(pastTime)
	if err != nil {
		t.Errorf("Expected no error for past feeding time, got %v", err)
	}
	if !feedingTime.Time().Equal(pastTime) {
		t.Errorf("Expected feeding time to be %v, got %v", pastTime, feedingTime.Time())
	}
}

func TestGender(t *testing.T) {
	// Test male gender
	if domain.Male != "male" {
		t.Errorf("Expected Male to be 'male', got %s", domain.Male)
	}

	// Test female gender
	if domain.Female != "female" {
		t.Errorf("Expected Female to be 'female', got %s", domain.Female)
	}
}

func TestHealthStatus(t *testing.T) {
	// Test healthy status
	if domain.Healthy != "healthy" {
		t.Errorf("Expected Healthy to be 'healthy', got %s", domain.Healthy)
	}

	// Test sick status
	if domain.Sick != "sick" {
		t.Errorf("Expected Sick to be 'sick', got %s", domain.Sick)
	}
}

func TestEnclosureType(t *testing.T) {
	// Test predator type
	if domain.Predator != "predator" {
		t.Errorf("Expected Predator to be 'predator', got %s", domain.Predator)
	}

	// Test herbivore type
	if domain.Herbivore != "herbivore" {
		t.Errorf("Expected Herbivore to be 'herbivore', got %s", domain.Herbivore)
	}

	// Test aviary type
	if domain.Aviary != "aviary" {
		t.Errorf("Expected Aviary to be 'aviary', got %s", domain.Aviary)
	}

	// Test aquarium type
	if domain.Aquarium != "aquarium" {
		t.Errorf("Expected Aquarium to be 'aquarium', got %s", domain.Aquarium)
	}

	// Test terrarium type
	if domain.Terrarium != "terrarium" {
		t.Errorf("Expected Terrarium to be 'terrarium', got %s", domain.Terrarium)
	}
}

func TestFoodType(t *testing.T) {
	// Test meat type
	if domain.Meat != "meat" {
		t.Errorf("Expected Meat to be 'meat', got %s", domain.Meat)
	}

	// Test vegetables type
	if domain.Vegetables != "vegetables" {
		t.Errorf("Expected Vegetables to be 'vegetables', got %s", domain.Vegetables)
	}

	// Test fruits type
	if domain.Fruits != "fruits" {
		t.Errorf("Expected Fruits to be 'fruits', got %s", domain.Fruits)
	}

	// Test insects type
	if domain.Insects != "insects" {
		t.Errorf("Expected Insects to be 'insects', got %s", domain.Insects)
	}

	// Test seeds type
	if domain.Seeds != "seeds" {
		t.Errorf("Expected Seeds to be 'seeds', got %s", domain.Seeds)
	}
}
