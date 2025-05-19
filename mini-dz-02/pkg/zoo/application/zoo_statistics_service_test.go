package application_test

import (
	"testing"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

func TestGetZooStatistics(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()

	service := application.NewZooStatisticsService(animalRepo, enclosureRepo)

	// Create test data
	// Animals
	lion := createTestAnimal("lion1", "enclosure1")
	lion.Species, _ = domain.NewSpecies("Lion")
	lion.Gender = domain.Male
	lion.HealthStatus = domain.Healthy

	tiger := createTestAnimal("tiger1", "enclosure1")
	tiger.Species, _ = domain.NewSpecies("Tiger")
	tiger.Gender = domain.Female
	tiger.HealthStatus = domain.Sick

	elephant := createTestAnimal("elephant1", "enclosure2")
	elephant.Species, _ = domain.NewSpecies("Elephant")
	elephant.Gender = domain.Male
	elephant.HealthStatus = domain.Healthy

	giraffe := createTestAnimal("giraffe1", "")
	giraffe.Species, _ = domain.NewSpecies("Giraffe")
	giraffe.Gender = domain.Female
	giraffe.HealthStatus = domain.Healthy

	// Enclosures
	enclosure1 := createTestEnclosure("enclosure1", 2)
	enclosure1.AddAnimal(lion.ID)
	enclosure1.AddAnimal(tiger.ID)

	enclosure2 := createTestEnclosure("enclosure2", 2)
	enclosure2.AddAnimal(elephant.ID)

	enclosure3 := createTestEnclosure("enclosure3", 2)

	// Save to repositories
	animalRepo.Save(lion)
	animalRepo.Save(tiger)
	animalRepo.Save(elephant)
	animalRepo.Save(giraffe)

	enclosureRepo.Save(enclosure1)
	enclosureRepo.Save(enclosure2)
	enclosureRepo.Save(enclosure3)

	// Test getting zoo statistics
	stats, err := service.GetZooStatistics()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check total counts
	if stats.TotalAnimals != 4 {
		t.Errorf("Expected 4 total animals, got %d", stats.TotalAnimals)
	}

	if stats.HealthyAnimals != 3 {
		t.Errorf("Expected 3 healthy animals, got %d", stats.HealthyAnimals)
	}

	if stats.SickAnimals != 1 {
		t.Errorf("Expected 1 sick animal, got %d", stats.SickAnimals)
	}

	if stats.TotalEnclosures != 3 {
		t.Errorf("Expected 3 total enclosures, got %d", stats.TotalEnclosures)
	}

	if stats.AvailableEnclosures != 1 {
		t.Errorf("Expected 1 available enclosure, got %d", stats.AvailableEnclosures)
	}

	if stats.FullEnclosures != 1 {
		t.Errorf("Expected 1 full enclosure, got %d", stats.FullEnclosures)
	}

	if stats.EmptyEnclosures != 1 {
		t.Errorf("Expected 1 empty enclosure, got %d", stats.EmptyEnclosures)
	}

	// Check animals by species
	if stats.AnimalsBySpecies["Lion"] != 1 {
		t.Errorf("Expected 1 Lion, got %d", stats.AnimalsBySpecies["Lion"])
	}

	if stats.AnimalsBySpecies["Tiger"] != 1 {
		t.Errorf("Expected 1 Tiger, got %d", stats.AnimalsBySpecies["Tiger"])
	}

	if stats.AnimalsBySpecies["Elephant"] != 1 {
		t.Errorf("Expected 1 Elephant, got %d", stats.AnimalsBySpecies["Elephant"])
	}

	if stats.AnimalsBySpecies["Giraffe"] != 1 {
		t.Errorf("Expected 1 Giraffe, got %d", stats.AnimalsBySpecies["Giraffe"])
	}

	// Check animals by gender
	if stats.AnimalsByGender[domain.Male] != 2 {
		t.Errorf("Expected 2 male animals, got %d", stats.AnimalsByGender[domain.Male])
	}

	if stats.AnimalsByGender[domain.Female] != 2 {
		t.Errorf("Expected 2 female animals, got %d", stats.AnimalsByGender[domain.Female])
	}

	// Check animals by enclosure
	if stats.AnimalsByEnclosure["enclosure1"] != 2 {
		t.Errorf("Expected 2 animals in enclosure1, got %d", stats.AnimalsByEnclosure["enclosure1"])
	}

	if stats.AnimalsByEnclosure["enclosure2"] != 1 {
		t.Errorf("Expected 1 animal in enclosure2, got %d", stats.AnimalsByEnclosure["enclosure2"])
	}
}

func TestGetAnimalCountBySpecies(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()

	service := application.NewZooStatisticsService(animalRepo, enclosureRepo)

	// Create test data
	lion1 := createTestAnimal("lion1", "enclosure1")
	lion1.Species, _ = domain.NewSpecies("Lion")

	lion2 := createTestAnimal("lion2", "enclosure1")
	lion2.Species, _ = domain.NewSpecies("Lion")

	tiger := createTestAnimal("tiger1", "enclosure2")
	tiger.Species, _ = domain.NewSpecies("Tiger")

	elephant := createTestAnimal("elephant1", "enclosure3")
	elephant.Species, _ = domain.NewSpecies("Elephant")

	// Save to repository
	animalRepo.Save(lion1)
	animalRepo.Save(lion2)
	animalRepo.Save(tiger)
	animalRepo.Save(elephant)

	// Test getting animal count by species
	animalsBySpecies, err := service.GetAnimalCountBySpecies()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if animalsBySpecies["Lion"] != 2 {
		t.Errorf("Expected 2 Lions, got %d", animalsBySpecies["Lion"])
	}

	if animalsBySpecies["Tiger"] != 1 {
		t.Errorf("Expected 1 Tiger, got %d", animalsBySpecies["Tiger"])
	}

	if animalsBySpecies["Elephant"] != 1 {
		t.Errorf("Expected 1 Elephant, got %d", animalsBySpecies["Elephant"])
	}
}

func TestGetEnclosureUtilization(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()

	service := application.NewZooStatisticsService(animalRepo, enclosureRepo)

	// Create test data
	enclosure1 := createTestEnclosure("enclosure1", 4)
	enclosure1.AddAnimal("animal1")
	enclosure1.AddAnimal("animal2")

	enclosure2 := createTestEnclosure("enclosure2", 2)
	enclosure2.AddAnimal("animal3")

	enclosure3 := createTestEnclosure("enclosure3", 3)

	// Save to repository
	enclosureRepo.Save(enclosure1)
	enclosureRepo.Save(enclosure2)
	enclosureRepo.Save(enclosure3)

	// Test getting enclosure utilization
	utilization, err := service.GetEnclosureUtilization()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check utilization percentages
	expectedUtilization1 := float64(2) / float64(4) // 2 animals in a capacity of 4
	if utilization["enclosure1"] != expectedUtilization1 {
		t.Errorf("Expected utilization of enclosure1 to be %f, got %f", expectedUtilization1, utilization["enclosure1"])
	}

	expectedUtilization2 := float64(1) / float64(2) // 1 animal in a capacity of 2
	if utilization["enclosure2"] != expectedUtilization2 {
		t.Errorf("Expected utilization of enclosure2 to be %f, got %f", expectedUtilization2, utilization["enclosure2"])
	}

	expectedUtilization3 := float64(0) / float64(3) // 0 animals in a capacity of 3
	if utilization["enclosure3"] != expectedUtilization3 {
		t.Errorf("Expected utilization of enclosure3 to be %f, got %f", expectedUtilization3, utilization["enclosure3"])
	}
}

func TestGetHealthStatusStatistics(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()

	service := application.NewZooStatisticsService(animalRepo, enclosureRepo)

	// Create test data
	animal1 := createTestAnimal("animal1", "enclosure1")
	animal1.HealthStatus = domain.Healthy

	animal2 := createTestAnimal("animal2", "enclosure1")
	animal2.HealthStatus = domain.Healthy

	animal3 := createTestAnimal("animal3", "enclosure2")
	animal3.HealthStatus = domain.Sick

	animal4 := createTestAnimal("animal4", "enclosure2")
	animal4.HealthStatus = domain.Sick

	animal5 := createTestAnimal("animal5", "enclosure3")
	animal5.HealthStatus = domain.Healthy

	// Save to repository
	animalRepo.Save(animal1)
	animalRepo.Save(animal2)
	animalRepo.Save(animal3)
	animalRepo.Save(animal4)
	animalRepo.Save(animal5)

	// Test getting health status statistics
	healthStats, err := service.GetHealthStatusStatistics()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if healthStats[domain.Healthy] != 3 {
		t.Errorf("Expected 3 healthy animals, got %d", healthStats[domain.Healthy])
	}

	if healthStats[domain.Sick] != 2 {
		t.Errorf("Expected 2 sick animals, got %d", healthStats[domain.Sick])
	}
}
