package application_test

import (
	"errors"
	"testing"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// Mock implementations for testing
type MockAnimalRepository struct {
	animals      map[string]*domain.Animal
	saveError    bool
	getByIDError map[string]bool
}

func NewMockAnimalRepository() *MockAnimalRepository {
	return &MockAnimalRepository{
		animals:      make(map[string]*domain.Animal),
		getByIDError: make(map[string]bool),
	}
}

func (m *MockAnimalRepository) SetSaveError(shouldError bool) {
	m.saveError = shouldError
}

func (m *MockAnimalRepository) SetGetByIDError(id string, shouldError bool) {
	m.getByIDError[id] = shouldError
}

func (m *MockAnimalRepository) GetByID(id string) (*domain.Animal, error) {
	if m.getByIDError[id] {
		return nil, errors.New("simulated error getting animal")
	}

	animal, exists := m.animals[id]
	if !exists {
		return nil, errors.New("animal not found")
	}
	return animal, nil
}

func (m *MockAnimalRepository) GetAll() ([]*domain.Animal, error) {
	animals := make([]*domain.Animal, 0, len(m.animals))
	for _, animal := range m.animals {
		animals = append(animals, animal)
	}
	return animals, nil
}

func (m *MockAnimalRepository) Save(animal *domain.Animal) error {
	if m.saveError {
		return errors.New("simulated error saving animal")
	}
	m.animals[animal.ID] = animal
	return nil
}

func (m *MockAnimalRepository) Delete(id string) error {
	delete(m.animals, id)
	return nil
}

func (m *MockAnimalRepository) GetByEnclosureID(enclosureID string) ([]*domain.Animal, error) {
	var result []*domain.Animal
	for _, animal := range m.animals {
		if animal.EnclosureID == enclosureID {
			result = append(result, animal)
		}
	}
	return result, nil
}

type MockEnclosureRepository struct {
	enclosures   map[string]*domain.Enclosure
	saveError    bool
	getByIDError map[string]bool
}

func NewMockEnclosureRepository() *MockEnclosureRepository {
	return &MockEnclosureRepository{
		enclosures:   make(map[string]*domain.Enclosure),
		getByIDError: make(map[string]bool),
	}
}

func (m *MockEnclosureRepository) SetSaveError(shouldError bool) {
	m.saveError = shouldError
}

func (m *MockEnclosureRepository) SetGetByIDError(id string, shouldError bool) {
	m.getByIDError[id] = shouldError
}

func (m *MockEnclosureRepository) GetByID(id string) (*domain.Enclosure, error) {
	if m.getByIDError[id] {
		return nil, errors.New("simulated error getting enclosure")
	}

	enclosure, exists := m.enclosures[id]
	if !exists {
		return nil, errors.New("enclosure not found")
	}
	return enclosure, nil
}

func (m *MockEnclosureRepository) GetAll() ([]*domain.Enclosure, error) {
	enclosures := make([]*domain.Enclosure, 0, len(m.enclosures))
	for _, enclosure := range m.enclosures {
		enclosures = append(enclosures, enclosure)
	}
	return enclosures, nil
}

func (m *MockEnclosureRepository) Save(enclosure *domain.Enclosure) error {
	if m.saveError {
		return errors.New("simulated error saving enclosure")
	}
	m.enclosures[enclosure.ID] = enclosure
	return nil
}

func (m *MockEnclosureRepository) Delete(id string) error {
	delete(m.enclosures, id)
	return nil
}

func (m *MockEnclosureRepository) GetByType(enclosureType domain.EnclosureType) ([]*domain.Enclosure, error) {
	var result []*domain.Enclosure
	for _, enclosure := range m.enclosures {
		if enclosure.Type == enclosureType {
			result = append(result, enclosure)
		}
	}
	return result, nil
}

func (m *MockEnclosureRepository) GetAvailable() ([]*domain.Enclosure, error) {
	var result []*domain.Enclosure
	for _, enclosure := range m.enclosures {
		if enclosure.HasSpace() {
			result = append(result, enclosure)
		}
	}
	return result, nil
}

type MockEventPublisher struct {
	publishedEvents []domain.Event
	publishError    bool
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		publishedEvents: make([]domain.Event, 0),
		publishError:    false,
	}
}

func (m *MockEventPublisher) SetPublishError(shouldError bool) {
	m.publishError = shouldError
}

func (m *MockEventPublisher) Publish(event domain.Event) error {
	if m.publishError {
		return errors.New("simulated error publishing event")
	}
	m.publishedEvents = append(m.publishedEvents, event)
	return nil
}

// Helper function to create a test animal
func createTestAnimal(id string, enclosureID string) *domain.Animal {
	species, _ := domain.NewSpecies("Lion")
	name, _ := domain.NewAnimalName("Leo")
	birthDate, _ := domain.NewBirthDate(domain.BirthDate{}.Time())
	animal := &domain.Animal{
		ID:           id,
		Species:      species,
		Name:         name,
		BirthDate:    birthDate,
		Gender:       domain.Male,
		FavoriteFood: domain.Meat,
		HealthStatus: domain.Healthy,
		EnclosureID:  enclosureID,
	}
	return animal
}

// Helper function to create a test enclosure
func createTestEnclosure(id string, maxCapacity int) *domain.Enclosure {
	size, _ := domain.NewEnclosureSize(maxCapacity * 2)
	capacity, _ := domain.NewCapacity(maxCapacity)
	enclosure := &domain.Enclosure{
		ID:               id,
		Type:             domain.Predator,
		Size:             size,
		MaxCapacity:      capacity,
		CurrentAnimalIDs: make([]string, 0),
		CleaningStatus:   true,
	}
	return enclosure
}

func TestTransferAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "enclosure1")
	fromEnclosure := createTestEnclosure("enclosure1", 5)
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Add animal to source enclosure
	fromEnclosure.AddAnimal(animal.ID)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(fromEnclosure)
	enclosureRepo.Save(toEnclosure)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if animal was moved
	updatedAnimal, _ := animalRepo.GetByID(animal.ID)
	if updatedAnimal.EnclosureID != toEnclosure.ID {
		t.Errorf("Expected animal to be in enclosure %s, got %s", toEnclosure.ID, updatedAnimal.EnclosureID)
	}

	// Check if animal was removed from source enclosure
	updatedFromEnclosure, _ := enclosureRepo.GetByID(fromEnclosure.ID)
	for _, id := range updatedFromEnclosure.CurrentAnimalIDs {
		if id == animal.ID {
			t.Errorf("Animal should have been removed from source enclosure")
		}
	}

	// Check if animal was added to target enclosure
	updatedToEnclosure, _ := enclosureRepo.GetByID(toEnclosure.ID)
	found := false
	for _, id := range updatedToEnclosure.CurrentAnimalIDs {
		if id == animal.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Animal should have been added to target enclosure")
	}

	// Check if event was published
	if len(eventPublisher.publishedEvents) != 1 {
		t.Errorf("Expected 1 event to be published, got %d", len(eventPublisher.publishedEvents))
	}
}

func TestTransferAnimal_AnimalNotFound(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Test transfer with non-existent animal
	err := service.TransferAnimal("nonexistent", "enclosure1")

	// Assertions
	if err == nil {
		t.Errorf("Expected error for non-existent animal, got nil")
	}
}

func TestTransferAnimal_EnclosureNotFound(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "")
	animalRepo.Save(animal)

	// Test transfer to non-existent enclosure
	err := service.TransferAnimal(animal.ID, "nonexistent")

	// Assertions
	if err == nil {
		t.Errorf("Expected error for non-existent enclosure, got nil")
	}
}

func TestTransferAnimal_EnclosureAtCapacity(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "")
	enclosure := createTestEnclosure("enclosure1", 1)

	// Fill the enclosure to capacity
	enclosure.AddAnimal("other-animal")

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(enclosure)

	// Test transfer to full enclosure
	err := service.TransferAnimal(animal.ID, enclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error for enclosure at capacity, got nil")
	}
}

func TestGetAvailableEnclosures(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	enclosure1 := createTestEnclosure("enclosure1", 2)
	enclosure2 := createTestEnclosure("enclosure2", 1)

	// Fill enclosure2 to capacity
	enclosure2.AddAnimal("animal1")

	// Save to repository
	enclosureRepo.Save(enclosure1)
	enclosureRepo.Save(enclosure2)

	// Test getting available enclosures
	availableEnclosures, err := service.GetAvailableEnclosures()

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(availableEnclosures) != 1 {
		t.Errorf("Expected 1 available enclosure, got %d", len(availableEnclosures))
	}

	if availableEnclosures[0].ID != enclosure1.ID {
		t.Errorf("Expected enclosure1 to be available, got %s", availableEnclosures[0].ID)
	}
}

func TestGetAnimalsByEnclosure(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	enclosureID := "enclosure1"
	animal1 := createTestAnimal("animal1", enclosureID)
	animal2 := createTestAnimal("animal2", enclosureID)
	animal3 := createTestAnimal("animal3", "enclosure2")

	// Save to repository
	animalRepo.Save(animal1)
	animalRepo.Save(animal2)
	animalRepo.Save(animal3)

	// Test getting animals by enclosure
	animals, err := service.GetAnimalsByEnclosure(enclosureID)

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(animals) != 2 {
		t.Errorf("Expected 2 animals in enclosure, got %d", len(animals))
	}

	// Check that we got the right animals
	animalIDs := make(map[string]bool)
	for _, animal := range animals {
		animalIDs[animal.ID] = true
	}

	if !animalIDs["animal1"] || !animalIDs["animal2"] {
		t.Errorf("Did not get the expected animals")
	}
}

func TestTransferAnimal_GetCurrentEnclosureError(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "enclosure1")
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Set up error for getting the current enclosure
	enclosureRepo.SetGetByIDError("enclosure1", true)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(toEnclosure)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error when getting current enclosure, got nil")
	}
}

func TestTransferAnimal_SaveCurrentEnclosureError(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "enclosure1")
	fromEnclosure := createTestEnclosure("enclosure1", 5)
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Add animal to source enclosure
	fromEnclosure.AddAnimal(animal.ID)

	// Set up error for saving the current enclosure
	enclosureRepo.SetSaveError(true)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(fromEnclosure)
	enclosureRepo.Save(toEnclosure)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error when saving current enclosure, got nil")
	}
}

func TestTransferAnimal_SaveAnimalError(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "")
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Set up error for saving the animal
	animalRepo.SetSaveError(true)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(toEnclosure)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error when saving animal, got nil")
	}
}

func TestTransferAnimal_SaveTargetEnclosureError(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "")
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(toEnclosure)

	// Set up error for saving the target enclosure
	// We need to do this after saving the enclosure initially
	enclosureRepo.SetSaveError(true)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error when saving target enclosure, got nil")
	}
}

func TestTransferAnimal_PublishEventError(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := NewMockEventPublisher()

	service := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	animal := createTestAnimal("animal1", "")
	toEnclosure := createTestEnclosure("enclosure2", 5)

	// Save to repositories
	animalRepo.Save(animal)
	enclosureRepo.Save(toEnclosure)

	// Set up error for publishing the event
	eventPublisher.SetPublishError(true)

	// Test transfer
	err := service.TransferAnimal(animal.ID, toEnclosure.ID)

	// Assertions
	if err == nil {
		t.Errorf("Expected error when publishing event, got nil")
	}
}
