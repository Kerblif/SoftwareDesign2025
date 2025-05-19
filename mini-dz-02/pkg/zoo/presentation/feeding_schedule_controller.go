package presentation

import (
	"encoding/json"
	"net/http"
	"time"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// FeedingScheduleController handles HTTP requests related to feeding schedules
type FeedingScheduleController struct {
	feedingScheduleRepository domain.FeedingScheduleRepository
	animalRepository          domain.AnimalRepository
	feedingService            *application.FeedingOrganizationService
}

// NewFeedingScheduleController creates a new FeedingScheduleController
func NewFeedingScheduleController(
	feedingScheduleRepository domain.FeedingScheduleRepository,
	animalRepository domain.AnimalRepository,
	feedingService *application.FeedingOrganizationService,
) *FeedingScheduleController {
	return &FeedingScheduleController{
		feedingScheduleRepository: feedingScheduleRepository,
		animalRepository:          animalRepository,
		feedingService:            feedingService,
	}
}

// FeedingScheduleResponse represents the response for a feeding schedule
type FeedingScheduleResponse struct {
	ID          string    `json:"id"`
	AnimalID    string    `json:"animalId"`
	FeedingTime time.Time `json:"feedingTime"`
	FoodType    string    `json:"foodType"`
	Completed   bool      `json:"completed"`
	IsDue       bool      `json:"isDue"`
}

// CreateFeedingScheduleRequest represents the request to create a feeding schedule
type CreateFeedingScheduleRequest struct {
	ID          string    `json:"id"`
	AnimalID    string    `json:"animalId"`
	FeedingTime time.Time `json:"feedingTime"`
	FoodType    string    `json:"foodType"`
}

// UpdateFeedingScheduleRequest represents the request to update a feeding schedule
type UpdateFeedingScheduleRequest struct {
	FeedingTime time.Time `json:"feedingTime"`
	FoodType    string    `json:"foodType"`
}

// RegisterRoutes registers the routes for the FeedingScheduleController
func (c *FeedingScheduleController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/feeding-schedules", c.GetAllFeedingSchedules)
	mux.HandleFunc("GET /api/feeding-schedules/{id}", c.GetFeedingSchedule)
	mux.HandleFunc("POST /api/feeding-schedules", c.CreateFeedingSchedule)
	mux.HandleFunc("DELETE /api/feeding-schedules/{id}", c.DeleteFeedingSchedule)
	mux.HandleFunc("PUT /api/feeding-schedules/{id}", c.UpdateFeedingSchedule)
	mux.HandleFunc("POST /api/feeding-schedules/{id}/complete", c.CompleteFeedingSchedule)
	mux.HandleFunc("GET /api/feeding-schedules/due", c.GetDueFeedingSchedules)
	mux.HandleFunc("GET /api/animals/{id}/feeding-schedules", c.GetFeedingSchedulesByAnimal)
}

// GetAllFeedingSchedules handles GET /api/feeding-schedules
func (c *FeedingScheduleController) GetAllFeedingSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := c.feedingService.GetAllFeedingSchedules()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []FeedingScheduleResponse
	for _, schedule := range schedules {
		response = append(response, FeedingScheduleResponse{
			ID:          schedule.ID,
			AnimalID:    schedule.AnimalID,
			FeedingTime: schedule.FeedingTime.Time(),
			FoodType:    string(schedule.FoodType),
			Completed:   schedule.IsCompleted(),
			IsDue:       schedule.IsDue(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFeedingSchedule handles GET /api/feeding-schedules/{id}
func (c *FeedingScheduleController) GetFeedingSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Feeding schedule ID is required", http.StatusBadRequest)
		return
	}

	schedule, err := c.feedingScheduleRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := FeedingScheduleResponse{
		ID:          schedule.ID,
		AnimalID:    schedule.AnimalID,
		FeedingTime: schedule.FeedingTime.Time(),
		FoodType:    string(schedule.FoodType),
		Completed:   schedule.IsCompleted(),
		IsDue:       schedule.IsDue(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateFeedingSchedule handles POST /api/feeding-schedules
func (c *FeedingScheduleController) CreateFeedingSchedule(w http.ResponseWriter, r *http.Request) {
	var request CreateFeedingScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the feeding schedule
	err := c.feedingService.CreateFeedingSchedule(
		request.ID,
		request.AnimalID,
		request.FeedingTime,
		domain.FoodType(request.FoodType),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the created schedule
	schedule, err := c.feedingScheduleRepository.GetByID(request.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := FeedingScheduleResponse{
		ID:          schedule.ID,
		AnimalID:    schedule.AnimalID,
		FeedingTime: schedule.FeedingTime.Time(),
		FoodType:    string(schedule.FoodType),
		Completed:   schedule.IsCompleted(),
		IsDue:       schedule.IsDue(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteFeedingSchedule handles DELETE /api/feeding-schedules/{id}
func (c *FeedingScheduleController) DeleteFeedingSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Feeding schedule ID is required", http.StatusBadRequest)
		return
	}

	// Check if the feeding schedule exists
	_, err := c.feedingScheduleRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete the feeding schedule
	if err := c.feedingScheduleRepository.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateFeedingSchedule handles PUT /api/feeding-schedules/{id}
func (c *FeedingScheduleController) UpdateFeedingSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Feeding schedule ID is required", http.StatusBadRequest)
		return
	}

	var request UpdateFeedingScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the feeding schedule
	err := c.feedingService.UpdateFeedingSchedule(
		id,
		request.FeedingTime,
		domain.FoodType(request.FoodType),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the updated schedule
	schedule, err := c.feedingScheduleRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := FeedingScheduleResponse{
		ID:          schedule.ID,
		AnimalID:    schedule.AnimalID,
		FeedingTime: schedule.FeedingTime.Time(),
		FoodType:    string(schedule.FoodType),
		Completed:   schedule.IsCompleted(),
		IsDue:       schedule.IsDue(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompleteFeedingSchedule handles POST /api/feeding-schedules/{id}/complete
func (c *FeedingScheduleController) CompleteFeedingSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Feeding schedule ID is required", http.StatusBadRequest)
		return
	}

	// Mark the feeding schedule as completed
	if err := c.feedingService.MarkFeedingCompleted(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetDueFeedingSchedules handles GET /api/feeding-schedules/due
func (c *FeedingScheduleController) GetDueFeedingSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := c.feedingService.GetDueFeedings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []FeedingScheduleResponse
	for _, schedule := range schedules {
		response = append(response, FeedingScheduleResponse{
			ID:          schedule.ID,
			AnimalID:    schedule.AnimalID,
			FeedingTime: schedule.FeedingTime.Time(),
			FoodType:    string(schedule.FoodType),
			Completed:   schedule.IsCompleted(),
			IsDue:       schedule.IsDue(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetFeedingSchedulesByAnimal handles GET /api/animals/{id}/feeding-schedules
func (c *FeedingScheduleController) GetFeedingSchedulesByAnimal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	// Check if the animal exists
	_, err := c.animalRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get feeding schedules for the animal
	schedules, err := c.feedingService.GetFeedingSchedulesByAnimal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []FeedingScheduleResponse
	for _, schedule := range schedules {
		response = append(response, FeedingScheduleResponse{
			ID:          schedule.ID,
			AnimalID:    schedule.AnimalID,
			FeedingTime: schedule.FeedingTime.Time(),
			FoodType:    string(schedule.FoodType),
			Completed:   schedule.IsCompleted(),
			IsDue:       schedule.IsDue(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}