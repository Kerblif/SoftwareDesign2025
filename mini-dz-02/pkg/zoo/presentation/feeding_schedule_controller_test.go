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

// MockFeedingScheduleRepository is a mock implementation of domain.FeedingScheduleRepository
type MockFeedingScheduleRepository struct {
	schedules    map[string]*domain.FeedingSchedule
	getByIDError bool
	getAllError  bool
	saveError    bool
	deleteError  bool
}

func NewMockFeedingScheduleRepository() *MockFeedingScheduleRepository {
	return &MockFeedingScheduleRepository{
		schedules: make(map[string]*domain.FeedingSchedule),
	}
}

func (m *MockFeedingScheduleRepository) GetByID(id string) (*domain.FeedingSchedule, error) {
	if m.getByIDError {
		return nil, errors.New("error getting feeding schedule")
	}
	schedule, exists := m.schedules[id]
	if !exists {
		return nil, errors.New("feeding schedule not found")
	}
	return schedule, nil
}

func (m *MockFeedingScheduleRepository) GetAll() ([]*domain.FeedingSchedule, error) {
	if m.getAllError {
		return nil, errors.New("error getting all feeding schedules")
	}
	schedules := make([]*domain.FeedingSchedule, 0, len(m.schedules))
	for _, schedule := range m.schedules {
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (m *MockFeedingScheduleRepository) Save(schedule *domain.FeedingSchedule) error {
	if m.saveError {
		return errors.New("error saving feeding schedule")
	}
	m.schedules[schedule.ID] = schedule
	return nil
}

func (m *MockFeedingScheduleRepository) Delete(id string) error {
	if m.deleteError {
		return errors.New("error deleting feeding schedule")
	}
	if _, exists := m.schedules[id]; !exists {
		return errors.New("feeding schedule not found")
	}
	delete(m.schedules, id)
	return nil
}

func (m *MockFeedingScheduleRepository) GetByAnimalID(animalID string) ([]*domain.FeedingSchedule, error) {
	var result []*domain.FeedingSchedule
	for _, schedule := range m.schedules {
		if schedule.AnimalID == animalID {
			result = append(result, schedule)
		}
	}
	return result, nil
}

func (m *MockFeedingScheduleRepository) GetDueSchedules() ([]*domain.FeedingSchedule, error) {
	var result []*domain.FeedingSchedule
	for _, schedule := range m.schedules {
		if schedule.IsDue() && !schedule.IsCompleted() {
			result = append(result, schedule)
		}
	}
	return result, nil
}

// Helper function to create a test feeding schedule
func createTestFeedingSchedule(id string, animalID string, completed bool) *domain.FeedingSchedule {
	feedingTime, _ := domain.NewFeedingTime(time.Now().Add(1 * time.Hour))
	schedule, _ := domain.NewFeedingSchedule(id, animalID, feedingTime, domain.Meat)
	if completed {
		schedule.MarkCompleted()
	}
	return schedule
}

// setupFeedingScheduleTest creates a controller and repositories for testing feeding schedule controller
func setupFeedingScheduleTest() (*presentation.FeedingScheduleController, *MockAnimalRepository, *MockFeedingScheduleRepository, *MockEventPublisher) {
	animalRepo := NewMockAnimalRepository()
	feedingScheduleRepo := NewMockFeedingScheduleRepository()
	eventPublisher := &MockEventPublisher{}

	// Create test data
	// Animals
	lion := createTestAnimal("animal-1")
	lion.Species, _ = domain.NewSpecies("Lion")
	lion.Gender = domain.Male
	lion.HealthStatus = domain.Healthy

	tiger := createTestAnimal("animal-2")
	tiger.Species, _ = domain.NewSpecies("Tiger")
	tiger.Gender = domain.Female
	tiger.HealthStatus = domain.Sick

	// Feeding schedules
	schedule1 := createTestFeedingSchedule("schedule-1", "animal-1", false)
	schedule2 := createTestFeedingSchedule("schedule-2", "animal-1", true)
	schedule3 := createTestFeedingSchedule("schedule-3", "animal-2", false)

	// Make schedule3 due
	pastFeedingTime, _ := domain.NewFeedingTime(time.Now().Add(-1 * time.Hour))
	schedule3.FeedingTime = pastFeedingTime

	// Save to repositories
	animalRepo.Save(lion)
	animalRepo.Save(tiger)

	feedingScheduleRepo.Save(schedule1)
	feedingScheduleRepo.Save(schedule2)
	feedingScheduleRepo.Save(schedule3)

	// Create service and controller
	feedingService := application.NewFeedingOrganizationService(animalRepo, feedingScheduleRepo, eventPublisher)
	controller := presentation.NewFeedingScheduleController(feedingScheduleRepo, animalRepo, feedingService)

	return controller, animalRepo, feedingScheduleRepo, eventPublisher
}

// TestGetAllFeedingSchedules tests the GetAllFeedingSchedules handler
func TestGetAllFeedingSchedules(t *testing.T) {
	// Setup
	controller, _, feedingScheduleRepo, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/feeding-schedules", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetAllFeedingSchedules(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 3 {
		t.Errorf("Expected 3 feeding schedules, got %d", len(response))
	}

	// Test error case
	feedingScheduleRepo.getAllError = true
	w = httptest.NewRecorder()
	controller.GetAllFeedingSchedules(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestGetFeedingSchedule tests the GetFeedingSchedule handler
func TestGetFeedingSchedule(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/feeding-schedules/schedule-1", nil)
	req = addPathParam(req, "id", "schedule-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetFeedingSchedule(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.ID != "schedule-1" {
		t.Errorf("Expected feeding schedule ID to be 'schedule-1', got %s", response.ID)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodGet, "/api/feeding-schedules/", nil)
	w = httptest.NewRecorder()
	controller.GetFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - feeding schedule not found
	req = httptest.NewRequest(http.MethodGet, "/api/feeding-schedules/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.GetFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}
}

// TestCreateFeedingSchedule tests the CreateFeedingSchedule handler
func TestCreateFeedingSchedule(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	createRequest := presentation.CreateFeedingScheduleRequest{
		ID:          "schedule-4",
		AnimalID:    "animal-1",
		FeedingTime: time.Now().Add(2 * time.Hour),
		FoodType:    "meat",
	}

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/feeding-schedules", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Call the handler
	controller.CreateFeedingSchedule(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Created, got %v", resp.Status)
	}

	var response presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.ID != "schedule-4" {
		t.Errorf("Expected feeding schedule ID to be 'schedule-4', got %s", response.ID)
	}

	// Test error case - invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/feeding-schedules", bytes.NewReader([]byte("invalid json")))
	w = httptest.NewRecorder()
	controller.CreateFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid animal ID
	createRequest.AnimalID = "non-existent"
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest(http.MethodPost, "/api/feeding-schedules", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.CreateFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

// TestDeleteFeedingSchedule tests the DeleteFeedingSchedule handler
func TestDeleteFeedingSchedule(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodDelete, "/api/feeding-schedules/schedule-1", nil)
	req = addPathParam(req, "id", "schedule-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.DeleteFeedingSchedule(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodDelete, "/api/feeding-schedules/", nil)
	w = httptest.NewRecorder()
	controller.DeleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - feeding schedule not found
	req = httptest.NewRequest(http.MethodDelete, "/api/feeding-schedules/non-existent", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.DeleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}

	// Test error case - delete error
	controller, _, feedingScheduleRepo, _ := setupFeedingScheduleTest()
	feedingScheduleRepo.deleteError = true
	req = httptest.NewRequest(http.MethodDelete, "/api/feeding-schedules/schedule-1", nil)
	req = addPathParam(req, "id", "schedule-1")
	w = httptest.NewRecorder()
	controller.DeleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status Internal Server Error, got %v", resp.Status)
	}
}

// TestUpdateFeedingSchedule tests the UpdateFeedingSchedule handler
func TestUpdateFeedingSchedule(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	updateRequest := presentation.UpdateFeedingScheduleRequest{
		FeedingTime: time.Now().Add(3 * time.Hour),
		FoodType:    "vegetables",
	}

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/api/feeding-schedules/schedule-1", bytes.NewReader(body))
	req = addPathParam(req, "id", "schedule-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.UpdateFeedingSchedule(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if response.ID != "schedule-1" {
		t.Errorf("Expected feeding schedule ID to be 'schedule-1', got %s", response.ID)
	}
	if response.FoodType != "vegetables" {
		t.Errorf("Expected food type to be 'vegetables', got %s", response.FoodType)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodPut, "/api/feeding-schedules/", bytes.NewReader(body))
	w = httptest.NewRecorder()
	controller.UpdateFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - invalid JSON
	req = httptest.NewRequest(http.MethodPut, "/api/feeding-schedules/schedule-1", bytes.NewReader([]byte("invalid json")))
	req = addPathParam(req, "id", "schedule-1")
	w = httptest.NewRecorder()
	controller.UpdateFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - feeding schedule not found
	req = httptest.NewRequest(http.MethodPut, "/api/feeding-schedules/non-existent", bytes.NewReader(body))
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.UpdateFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

// TestCompleteFeedingSchedule tests the CompleteFeedingSchedule handler
func TestCompleteFeedingSchedule(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/api/feeding-schedules/schedule-1/complete", nil)
	req = addPathParam(req, "id", "schedule-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.CompleteFeedingSchedule(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status No Content, got %v", resp.Status)
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodPost, "/api/feeding-schedules//complete", nil)
	w = httptest.NewRecorder()
	controller.CompleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - feeding schedule not found
	req = httptest.NewRequest(http.MethodPost, "/api/feeding-schedules/non-existent/complete", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.CompleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - already completed
	req = httptest.NewRequest(http.MethodPost, "/api/feeding-schedules/schedule-2/complete", nil)
	req = addPathParam(req, "id", "schedule-2")
	w = httptest.NewRecorder()
	controller.CompleteFeedingSchedule(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

// TestGetDueFeedingSchedules tests the GetDueFeedingSchedules handler
func TestGetDueFeedingSchedules(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/feeding-schedules/due", nil)
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetDueFeedingSchedules(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 1 {
		t.Errorf("Expected 1 due feeding schedule, got %d", len(response))
	}
	if len(response) > 0 && response[0].ID != "schedule-3" {
		t.Errorf("Expected due feeding schedule ID to be 'schedule-3', got %s", response[0].ID)
	}
}

// TestGetFeedingSchedulesByAnimal tests the GetFeedingSchedulesByAnimal handler
func TestGetFeedingSchedulesByAnimal(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/api/animals/animal-1/feeding-schedules", nil)
	req = addPathParam(req, "id", "animal-1")
	w := httptest.NewRecorder()

	// Call the handler
	controller.GetFeedingSchedulesByAnimal(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response []presentation.FeedingScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	// Verify the response data
	if len(response) != 2 {
		t.Errorf("Expected 2 feeding schedules for animal-1, got %d", len(response))
	}

	// Test error case - missing ID
	req = httptest.NewRequest(http.MethodGet, "/api/animals//feeding-schedules", nil)
	w = httptest.NewRecorder()
	controller.GetFeedingSchedulesByAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}

	// Test error case - animal not found
	req = httptest.NewRequest(http.MethodGet, "/api/animals/non-existent/feeding-schedules", nil)
	req = addPathParam(req, "id", "non-existent")
	w = httptest.NewRecorder()
	controller.GetFeedingSchedulesByAnimal(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %v", resp.Status)
	}
}

// TestRegisterFeedingScheduleRoutes tests the RegisterRoutes method
func TestRegisterFeedingScheduleRoutes(t *testing.T) {
	// Setup
	controller, _, _, _ := setupFeedingScheduleTest()
	mux := http.NewServeMux()

	// Call the method
	controller.RegisterRoutes(mux)

	// We can't easily test the routes directly, but we can at least verify that the method doesn't panic
	// In a real test, we might use a more sophisticated approach to verify the routes
}
