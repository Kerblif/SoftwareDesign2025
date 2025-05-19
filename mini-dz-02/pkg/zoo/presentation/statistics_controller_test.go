package presentation_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
	"mini-dz-02/pkg/zoo/presentation"
)

// setupStatisticsTest creates a controller and repositories for testing statistics controller
func setupStatisticsTest() (*presentation.StatisticsController, *MockAnimalRepository, *MockEnclosureRepository) {
	animalRepo := NewMockAnimalRepository()
	enclosureRepo := NewMockEnclosureRepository()

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

	// Create service and controller
	service := application.NewZooStatisticsService(animalRepo, enclosureRepo)
	controller := presentation.NewStatisticsController(service)

	return controller, animalRepo, enclosureRepo
}

// TestGetZooStatistics tests the GetZooStatistics handler
func TestGetZooStatistics(t *testing.T) {
	// Setup
	controller, animalRepo, _ := setupStatisticsTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/statistics", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetZooStatistics(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response presentation.ZooStatisticsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.TotalAnimals != 3 {
		t.Errorf("Expected TotalAnimals to be 3, got %d", response.TotalAnimals)
	}
	if response.HealthyAnimals != 2 {
		t.Errorf("Expected HealthyAnimals to be 2, got %d", response.HealthyAnimals)
	}
	if response.SickAnimals != 1 {
		t.Errorf("Expected SickAnimals to be 1, got %d", response.SickAnimals)
	}

	// Test error case - animal repository error
	controller, animalRepo, _ = setupStatisticsTest()
	animalRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetZooStatistics(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestGetAnimalCountBySpecies tests the GetAnimalCountBySpecies handler
func TestGetAnimalCountBySpecies(t *testing.T) {
	// Setup
	controller, animalRepo, _ := setupStatisticsTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/statistics/species", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAnimalCountBySpecies(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 3 {
		t.Errorf("Expected 3 species, got %d", len(response))
	}
	if response["Lion"] != 1 {
		t.Errorf("Expected 1 Lion, got %d", response["Lion"])
	}
	if response["Tiger"] != 1 {
		t.Errorf("Expected 1 Tiger, got %d", response["Tiger"])
	}
	if response["Elephant"] != 1 {
		t.Errorf("Expected 1 Elephant, got %d", response["Elephant"])
	}

	// Test error case
	controller, animalRepo, _ = setupStatisticsTest()
	animalRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetAnimalCountBySpecies(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestGetEnclosureUtilization tests the GetEnclosureUtilization handler
func TestGetEnclosureUtilization(t *testing.T) {
	// Setup
	controller, _, enclosureRepo := setupStatisticsTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/statistics/enclosure-utilization", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetEnclosureUtilization(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 3 {
		t.Errorf("Expected 3 enclosures, got %d", len(response))
	}

	// Test error case
	controller, _, enclosureRepo = setupStatisticsTest()
	enclosureRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetEnclosureUtilization(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestGetHealthStatusStatistics tests the GetHealthStatusStatistics handler
func TestGetHealthStatusStatistics(t *testing.T) {
	// Setup
	controller, animalRepo, _ := setupStatisticsTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/statistics/health", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetHealthStatusStatistics(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 2 {
		t.Errorf("Expected 2 health statuses, got %d", len(response))
	}
	if response["healthy"] != 2 {
		t.Errorf("Expected 2 healthy animals, got %d", response["healthy"])
	}
	if response["sick"] != 1 {
		t.Errorf("Expected 1 sick animal, got %d", response["sick"])
	}

	// Test error case
	controller, animalRepo, _ = setupStatisticsTest()
	animalRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetHealthStatusStatistics(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestRegisterStatisticsRoutes tests the RegisterRoutes method of the StatisticsController
func TestRegisterStatisticsRoutes(t *testing.T) {
	// Setup
	controller, _, _ := setupStatisticsTest()
	mux := http.NewServeMux()

	// Call the method
	controller.RegisterRoutes(mux)

	// We can't easily test the routes directly, but we can at least verify that the method doesn't panic
	// In a real test, we might use a more sophisticated approach to verify the routes
}
