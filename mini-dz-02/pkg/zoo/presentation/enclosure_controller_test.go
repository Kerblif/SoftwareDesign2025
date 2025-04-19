package presentation_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
	"mini-dz-02/pkg/zoo/presentation"
)

// setupEnclosureTest creates a controller and repositories for testing enclosure controller
func setupEnclosureTest() (*presentation.EnclosureController, *MockAnimalRepository, *MockEnclosureRepository) {
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()
	eventPublisher := &MockEventPublisher{}
	transferService := application.NewAnimalTransferService(animalRepo, enclosureRepo, eventPublisher)

	// Create test data
	// Animals
	lion := createTestAnimal("animal-1")
	lion.Species, _ = domain.NewSpecies("Lion")
	lion.Gender = domain.Male
	lion.HealthStatus = domain.Healthy
	lion.EnclosureID = "enclosure-1"

	tiger := createTestAnimal("animal-2")
	tiger.Species, _ = domain.NewSpecies("Tiger")
	tiger.Gender = domain.Female
	tiger.HealthStatus = domain.Sick
	tiger.EnclosureID = "enclosure-1"

	elephant := createTestAnimal("animal-3")
	elephant.Species, _ = domain.NewSpecies("Elephant")
	elephant.Gender = domain.Male
	elephant.HealthStatus = domain.Healthy
	elephant.EnclosureID = "enclosure-2"

	// Enclosures
	enclosure1 := createTestEnclosure("enclosure-1")
	enclosure1.AddAnimal(lion.ID)
	enclosure1.AddAnimal(tiger.ID)

	enclosure2 := createTestEnclosure("enclosure-2")
	enclosure2.AddAnimal(elephant.ID)

	enclosure3 := createTestEnclosure("enclosure-3")

	// Save to repositories
	animalRepo.Save(lion)
	animalRepo.Save(tiger)
	animalRepo.Save(elephant)

	enclosureRepo.Save(enclosure1)
	enclosureRepo.Save(enclosure2)
	enclosureRepo.Save(enclosure3)

	// Create controller
	controller := presentation.NewEnclosureController(enclosureRepo, animalRepo, transferService)

	return controller, animalRepo, enclosureRepo
}

// TestGetAllEnclosures tests the GetAllEnclosures handler
func TestGetAllEnclosures(t *testing.T) {
	// Setup
	controller, _, enclosureRepo := setupEnclosureTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/enclosures", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAllEnclosures(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.EnclosureResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 3 {
		t.Errorf("Expected 3 enclosures, got %d", len(response))
	}

	// Test error case
	enclosureRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetAllEnclosures(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestGetEnclosure tests the GetEnclosure handler
func TestGetEnclosure(t *testing.T) {
	// Setup
	controller, _, _ := setupEnclosureTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/enclosures/enclosure-1", nil)
	req = addPathParam(req, "id", "enclosure-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetEnclosure(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response presentation.EnclosureResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.ID != "enclosure-1" {
		t.Errorf("Expected enclosure ID to be 'enclosure-1', got %s", response.ID)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodGet, "/api/enclosures/", nil)
	w = httptest.NewRecorder()
	controller.GetEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - enclosure not found
	req = httptest.NewRequest(http.MethodGet, "/api/enclosures/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.GetEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}
}

// TestCreateEnclosure tests the CreateEnclosure handler
func TestCreateEnclosure(t *testing.T) {
	// Setup
	controller, _, _ := setupEnclosureTest()

	// Create a request
	createRequest := presentation.CreateEnclosureRequest{
		ID:          "enclosure-4",
		Type:        "predator",
		Size:        10,
		MaxCapacity: 5,
	}

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the handler
	controller.CreateEnclosure(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	var response presentation.EnclosureResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.ID != "enclosure-4" {
		t.Errorf("Expected enclosure ID to be 'enclosure-4', got %s", response.ID)
	}

	// Test error case - invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader([]byte("invalid json")))
	w = httptest.NewRecorder()
	controller.CreateEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid size
	createRequest.Size = -1
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.CreateEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid capacity
	createRequest.Size = 10
	createRequest.MaxCapacity = -1
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.CreateEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - size < capacity
	createRequest.Size = 5
	createRequest.MaxCapacity = 10
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.CreateEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - save error
	createRequest.Size = 10
	createRequest.MaxCapacity = 5
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller, _, enclosureRepo := setupEnclosureTest()
	enclosureRepo.saveError = true
	controller.CreateEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestDeleteEnclosure tests the DeleteEnclosure handler
func TestDeleteEnclosure(t *testing.T) {
	// Setup
	controller, _, _ := setupEnclosureTest()

	// Create a request for an empty enclosure
	req := httptest.NewRequest(http.MethodDelete, "/api/enclosures/enclosure-3", nil)
	req = addPathParam(req, "id", "enclosure-3")
	w := httptest.NewRecorder()

	// Call the handler
	controller.DeleteEnclosure(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodDelete, "/api/enclosures/", nil)
	w = httptest.NewRecorder()
	controller.DeleteEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - enclosure not found
	req = httptest.NewRequest(http.MethodDelete, "/api/enclosures/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.DeleteEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}

	// Test error case - enclosure not empty
	req = httptest.NewRequest(http.MethodDelete, "/api/enclosures/enclosure-1", nil)
	req = addPathParam(req, "id", "enclosure-1")
	w = httptest.NewRecorder()
	controller.DeleteEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

// TestGetAnimalsInEnclosure tests the GetAnimalsInEnclosure handler
func TestGetAnimalsInEnclosure(t *testing.T) {
	// Setup
	controller, _, _ := setupEnclosureTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/enclosures/enclosure-1/animals", nil)
	req = addPathParam(req, "id", "enclosure-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAnimalsInEnclosure(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.AnimalResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 2 {
		t.Errorf("Expected 2 animals, got %d", len(response))
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodGet, "/api/enclosures//animals", nil)
	w = httptest.NewRecorder()
	controller.GetAnimalsInEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - enclosure not found
	req = httptest.NewRequest(http.MethodGet, "/api/enclosures/non-existent/animals", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.GetAnimalsInEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}
}

// TestCleanEnclosure tests the CleanEnclosure handler
func TestCleanEnclosure(t *testing.T) {
	// Setup
	controller, _, enclosureRepo := setupEnclosureTest()

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/api/enclosures/enclosure-1/clean", nil)
	req = addPathParam(req, "id", "enclosure-1")
	w := httptest.NewRecorder()

	// Make the enclosure dirty first
	enclosure, _ := enclosureRepo.GetByID("enclosure-1")
	enclosure.CleaningStatus = false
	enclosureRepo.Save(enclosure)

	// Call the handler
	controller.CleanEnclosure(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures//clean", nil)
	w = httptest.NewRecorder()
	controller.CleanEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - enclosure not found
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures/non-existent/clean", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.CleanEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}

	// Test error case - already clean
	req = httptest.NewRequest(http.MethodPost, "/api/enclosures/enclosure-1/clean", nil)
	req = addPathParam(req, "id", "enclosure-1")
	w = httptest.NewRecorder()
	controller.CleanEnclosure(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

// TestRegisterEnclosureRoutes tests the RegisterRoutes method
func TestRegisterEnclosureRoutes(t *testing.T) {
	// Setup
	controller, _, _ := setupEnclosureTest()
	mux := http.NewServeMux()

	// Call the method
	controller.RegisterRoutes(mux)

	// We can't easily test the routes directly, but we can at least verify that the method doesn't panic
	// In a real test, we might use a more sophisticated approach to verify the routes
}
