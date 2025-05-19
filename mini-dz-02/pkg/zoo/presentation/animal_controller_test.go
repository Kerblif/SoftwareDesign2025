package presentation_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
	"mini-dz-02/pkg/zoo/presentation"
)

// MockAnimalRepository is a mock implementation of domain.AnimalRepository
type MockAnimalRepository struct {
	animals      map[string]*domain.Animal
	getByIDError bool
	getAllError  bool
	saveError    bool
	deleteError  bool
}

func NewMockAnimalRepository() *MockAnimalRepository {
	return &MockAnimalRepository{
		animals: make(map[string]*domain.Animal),
	}
}

func (m *MockAnimalRepository) GetByID(id string) (*domain.Animal, error) {
	if m.getByIDError {
		return nil, errors.New("error getting animal")
	}
	animal, exists := m.animals[id]
	if !exists {
		return nil, errors.New("animal not found")
	}
	return animal, nil
}

func (m *MockAnimalRepository) GetAll() ([]*domain.Animal, error) {
	if m.getAllError {
		return nil, errors.New("error getting all animals")
	}
	animals := make([]*domain.Animal, 0, len(m.animals))
	for _, animal := range m.animals {
		animals = append(animals, animal)
	}
	return animals, nil
}

func (m *MockAnimalRepository) Save(animal *domain.Animal) error {
	if m.saveError {
		return errors.New("error saving animal")
	}
	m.animals[animal.ID] = animal
	return nil
}

func (m *MockAnimalRepository) Delete(id string) error {
	if m.deleteError {
		return errors.New("error deleting animal")
	}
	if _, exists := m.animals[id]; !exists {
		return errors.New("animal not found")
	}
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

// MockEnclosureRepository is a mock implementation of domain.EnclosureRepository
type MockEnclosureRepository struct {
	enclosures   map[string]*domain.Enclosure
	getByIDError bool
	saveError    bool
	getAllError  bool
}

func NewMockEnclosureRepository() *MockEnclosureRepository {
	return &MockEnclosureRepository{
		enclosures: make(map[string]*domain.Enclosure),
	}
}

func (m *MockEnclosureRepository) GetByID(id string) (*domain.Enclosure, error) {
	if m.getByIDError {
		return nil, errors.New("error getting enclosure")
	}
	enclosure, exists := m.enclosures[id]
	if !exists {
		return nil, errors.New("enclosure not found")
	}
	return enclosure, nil
}

func (m *MockEnclosureRepository) GetAll() ([]*domain.Enclosure, error) {
	if m.getAllError {
		return nil, errors.New("error getting all enclosures")
	}
	enclosures := make([]*domain.Enclosure, 0, len(m.enclosures))
	for _, enclosure := range m.enclosures {
		enclosures = append(enclosures, enclosure)
	}
	return enclosures, nil
}

func (m *MockEnclosureRepository) Save(enclosure *domain.Enclosure) error {
	if m.saveError {
		return errors.New("error saving enclosure")
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

// MockTransferService is a mock implementation of the AnimalTransferService
type MockTransferService struct {
	transferError bool
}

func NewMockTransferService() *application.AnimalTransferService {
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := &MockEventPublisher{}
	return application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)
}

func NewMockTransferServiceWithRepos(animalRepo *MockAnimalRepository, enclosureRepo *MockEnclosureRepository) *application.AnimalTransferService {
	eventPublisher := &MockEventPublisher{}
	return application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)
}

// MockEventPublisher is a mock implementation of domain.EventPublisher
type MockEventPublisher struct {
	publishError bool
}

func (m *MockEventPublisher) Publish(event domain.Event) error {
	if m.publishError {
		return errors.New("error publishing event")
	}
	return nil
}

// Helper function to create a test animal
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

// Helper function to create a test enclosure
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

func TestGetAllAnimals(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	transferService := NewMockTransferService()

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	// Add some test animals
	animal1 := createTestAnimal("animal-1")
	animal2 := createTestAnimal("animal-2")
	animalRepo.Save(animal1)
	animalRepo.Save(animal2)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/animals", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAllAnimals(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.AnimalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 animals, got %d", len(response))
	}

	// Test error case
	animalRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetAllAnimals(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

func TestGetAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	transferService := NewMockTransferService()

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	// Add a test animal
	animal := createTestAnimal("animal-1")
	animalRepo.Save(animal)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/animals/animal-1", nil)
	req = addPathParam(req, "id", "animal-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response presentation.AnimalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.ID != "animal-1" {
		t.Errorf("Expected animal ID to be 'animal-1', got %s", response.ID)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodGet, "/api/animals/", nil)
	w = httptest.NewRecorder()
	controller.GetAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - animal not found
	req = httptest.NewRequest(http.MethodGet, "/api/animals/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.GetAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}
}

func TestCreateAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	transferService := NewMockTransferService()

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	// Create a request
	createRequest := presentation.CreateAnimalRequest{
		ID:           "animal-1",
		Species:      "Lion",
		Name:         "Leo",
		BirthDate:    time.Now().AddDate(-5, 0, 0),
		Gender:       "male",
		FavoriteFood: "meat",
		HealthStatus: "healthy",
	}

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/animals", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the handler
	controller.CreateAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	var response presentation.AnimalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if response.ID != "animal-1" {
		t.Errorf("Expected animal ID to be 'animal-1', got %s", response.ID)
	}

	// Test error case - invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/animals", bytes.NewReader([]byte("invalid json")))
	w = httptest.NewRecorder()
	controller.CreateAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid species
	createRequest.Species = ""
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/animals", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.CreateAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - save error
	createRequest.Species = "Lion"
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/animals", bytes.NewReader(body))
	w = httptest.NewRecorder()
	animalRepo.saveError = true
	controller.CreateAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

func TestDeleteAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	transferService := NewMockTransferService()

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	// Add a test animal
	animal := createTestAnimal("animal-1")
	animalRepo.Save(animal)

	// Create a request
	req := httptest.NewRequest(http.MethodDelete, "/api/animals/animal-1", nil)
	req = addPathParam(req, "id", "animal-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.DeleteAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodDelete, "/api/animals/", nil)
	w = httptest.NewRecorder()
	controller.DeleteAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - animal not found
	req = httptest.NewRequest(http.MethodDelete, "/api/animals/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.DeleteAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}

	// Test error case - delete error
	animal = createTestAnimal("animal-2")
	animalRepo.Save(animal)
	req = httptest.NewRequest(http.MethodDelete, "/api/animals/animal-2", nil)
	req = addPathParam(req, "id", "animal-2")
	w = httptest.NewRecorder()
	animalRepo.deleteError = true
	controller.DeleteAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}

	// Test case - animal in enclosure
	animalRepo.deleteError = false
	animal = createTestAnimal("animal-3")
	animal.EnclosureID = "enclosure-1"
	animalRepo.Save(animal)

	enclosure := createTestEnclosure("enclosure-1")
	enclosure.CurrentAnimalIDs = append(enclosure.CurrentAnimalIDs, "animal-3")
	enclosureRepo.Save(enclosure)

	req = httptest.NewRequest(http.MethodDelete, "/api/animals/animal-3", nil)
	req = addPathParam(req, "id", "animal-3")
	w = httptest.NewRecorder()
	controller.DeleteAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - enclosure not found
	animal = createTestAnimal("animal-4")
	animal.EnclosureID = "non-existent-enclosure"
	animalRepo.Save(animal)

	req = httptest.NewRequest(http.MethodDelete, "/api/animals/animal-4", nil)
	req = addPathParam(req, "id", "animal-4")
	w = httptest.NewRecorder()
	controller.DeleteAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

func TestTransferAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := &MockEventPublisher{}
	transferService := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	animal := createTestAnimal("animal-1")
	animalRepo.Save(animal)

	enclosure := createTestEnclosure("enclosure-1")
	enclosureRepo.Save(enclosure)

	// Create a request
	transferRequest := presentation.TransferAnimalRequest{
		EnclosureID: "enclosure-1",
	}

	body, _ := json.Marshal(transferRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/animals/animal-1/transfer", bytes.NewReader(body))
	req = addPathParam(req, "id", "animal-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.TransferAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodPost, "/api/animals//transfer", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.TransferAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/animals/animal-1/transfer", bytes.NewReader([]byte("invalid json")))
	req = addPathParam(req, "id", "animal-1")
	w = httptest.NewRecorder()
	controller.TransferAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - missing enclosure ID
	transferRequest.EnclosureID = ""
	body, _ = json.Marshal(transferRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/animals/animal-1/transfer", bytes.NewReader(body))
	req = addPathParam(req, "id", "animal-1")
	w = httptest.NewRecorder()
	controller.TransferAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Note: We can't easily test the transfer error case with the real service
	// In a real test, we would use a mock that allows setting error conditions
}

func TestTreatAnimal(t *testing.T) {
	// Setup
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	transferService := NewMockTransferService()

	controller := presentation.NewAnimalController(animalRepo, enclosureRepo, transferService)

	// Add a test animal
	animal := createTestAnimal("animal-1")
	animal.SetHealthStatus(domain.Sick)
	animalRepo.Save(animal)

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/api/animals/animal-1/treat", nil)
	req = addPathParam(req, "id", "animal-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.TreatAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodPost, "/api/animals//treat", nil)
	w = httptest.NewRecorder()
	controller.TreatAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - animal not found
	req = httptest.NewRequest(http.MethodPost, "/api/animals/non-existent/treat", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.TreatAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}

	// Test error case - already healthy
	animal = createTestAnimal("animal-2")
	animalRepo.Save(animal)
	req = httptest.NewRequest(http.MethodPost, "/api/animals/animal-2/treat", nil)
	req = addPathParam(req, "id", "animal-2")
	w = httptest.NewRecorder()
	controller.TreatAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - save error
	animal = createTestAnimal("animal-3")
	animal.SetHealthStatus(domain.Sick)
	animalRepo.Save(animal)
	req = httptest.NewRequest(http.MethodPost, "/api/animals/animal-3/treat", nil)
	req = addPathParam(req, "id", "animal-3")
	w = httptest.NewRecorder()
	animalRepo.saveError = true
	controller.TreatAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// Helper function to add path parameters to a request
func addPathParam(req *http.Request, key, value string) *http.Request {
	ctx := req.Context()
	req = req.WithContext(ctx)
	req.SetPathValue(key, value)
	return req
}
