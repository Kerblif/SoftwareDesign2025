package presentation

import (
	"encoding/json"
	"net/http"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// EnclosureController handles HTTP requests related to enclosures
type EnclosureController struct {
	enclosureRepository domain.EnclosureRepository
	animalRepository    domain.AnimalRepository
	transferService     *application.AnimalTransferService
}

// NewEnclosureController creates a new EnclosureController
func NewEnclosureController(
	enclosureRepository domain.EnclosureRepository,
	animalRepository domain.AnimalRepository,
	transferService *application.AnimalTransferService,
) *EnclosureController {
	return &EnclosureController{
		enclosureRepository: enclosureRepository,
		animalRepository:    animalRepository,
		transferService:     transferService,
	}
}

// EnclosureResponse represents the response for an enclosure
type EnclosureResponse struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	Size             int      `json:"size"`
	MaxCapacity      int      `json:"maxCapacity"`
	CurrentAnimalIDs []string `json:"currentAnimalIds"`
	IsClean          bool     `json:"isClean"`
}

// CreateEnclosureRequest represents the request to create an enclosure
type CreateEnclosureRequest struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	MaxCapacity int    `json:"maxCapacity"`
}

// RegisterRoutes registers the routes for the EnclosureController
func (c *EnclosureController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/enclosures", c.GetAllEnclosures)
	mux.HandleFunc("GET /api/enclosures/{id}", c.GetEnclosure)
	mux.HandleFunc("POST /api/enclosures", c.CreateEnclosure)
	mux.HandleFunc("DELETE /api/enclosures/{id}", c.DeleteEnclosure)
	mux.HandleFunc("GET /api/enclosures/{id}/animals", c.GetAnimalsInEnclosure)
	mux.HandleFunc("POST /api/enclosures/{id}/clean", c.CleanEnclosure)
}

// GetAllEnclosures handles GET /api/enclosures
func (c *EnclosureController) GetAllEnclosures(w http.ResponseWriter, r *http.Request) {
	enclosures, err := c.enclosureRepository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []EnclosureResponse
	for _, enclosure := range enclosures {
		response = append(response, EnclosureResponse{
			ID:               enclosure.ID,
			Type:             string(enclosure.Type),
			Size:             enclosure.Size.Value(),
			MaxCapacity:      enclosure.MaxCapacity.Value(),
			CurrentAnimalIDs: enclosure.CurrentAnimalIDs,
			IsClean:          enclosure.IsClean(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEnclosure handles GET /api/enclosures/{id}
func (c *EnclosureController) GetEnclosure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Enclosure ID is required", http.StatusBadRequest)
		return
	}

	enclosure, err := c.enclosureRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := EnclosureResponse{
		ID:               enclosure.ID,
		Type:             string(enclosure.Type),
		Size:             enclosure.Size.Value(),
		MaxCapacity:      enclosure.MaxCapacity.Value(),
		CurrentAnimalIDs: enclosure.CurrentAnimalIDs,
		IsClean:          enclosure.IsClean(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateEnclosure handles POST /api/enclosures
func (c *EnclosureController) CreateEnclosure(w http.ResponseWriter, r *http.Request) {
	var request CreateEnclosureRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create value objects
	size, err := domain.NewEnclosureSize(request.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	capacity, err := domain.NewCapacity(request.MaxCapacity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the enclosure
	enclosure, err := domain.NewEnclosure(
		request.ID,
		domain.EnclosureType(request.Type),
		size,
		capacity,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the enclosure
	if err := c.enclosureRepository.Save(enclosure); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := EnclosureResponse{
		ID:               enclosure.ID,
		Type:             string(enclosure.Type),
		Size:             enclosure.Size.Value(),
		MaxCapacity:      enclosure.MaxCapacity.Value(),
		CurrentAnimalIDs: enclosure.CurrentAnimalIDs,
		IsClean:          enclosure.IsClean(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteEnclosure handles DELETE /api/enclosures/{id}
func (c *EnclosureController) DeleteEnclosure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Enclosure ID is required", http.StatusBadRequest)
		return
	}

	// Check if the enclosure exists
	enclosure, err := c.enclosureRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if the enclosure is empty
	if !enclosure.IsEmpty() {
		http.Error(w, "Cannot delete enclosure that contains animals", http.StatusBadRequest)
		return
	}

	// Delete the enclosure
	if err := c.enclosureRepository.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAnimalsInEnclosure handles GET /api/enclosures/{id}/animals
func (c *EnclosureController) GetAnimalsInEnclosure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Enclosure ID is required", http.StatusBadRequest)
		return
	}

	// Check if the enclosure exists
	_, err := c.enclosureRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get animals in the enclosure
	animals, err := c.transferService.GetAnimalsByEnclosure(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []AnimalResponse
	for _, animal := range animals {
		response = append(response, AnimalResponse{
			ID:           animal.ID,
			Species:      animal.Species.String(),
			Name:         animal.Name.String(),
			BirthDate:    animal.BirthDate.Time(),
			Gender:       string(animal.Gender),
			FavoriteFood: string(animal.FavoriteFood),
			HealthStatus: string(animal.HealthStatus),
			EnclosureID:  animal.EnclosureID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CleanEnclosure handles POST /api/enclosures/{id}/clean
func (c *EnclosureController) CleanEnclosure(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Enclosure ID is required", http.StatusBadRequest)
		return
	}

	// Get the enclosure
	enclosure, err := c.enclosureRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Clean the enclosure
	if err := enclosure.Clean(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the enclosure
	if err := c.enclosureRepository.Save(enclosure); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}