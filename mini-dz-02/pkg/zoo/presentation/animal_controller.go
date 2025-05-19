package presentation

import (
	"encoding/json"
	"net/http"
	"time"

	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// AnimalController handles HTTP requests related to animals
type AnimalController struct {
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
	transferService     *application.AnimalTransferService
}

// NewAnimalController creates a new AnimalController
func NewAnimalController(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
	transferService *application.AnimalTransferService,
) *AnimalController {
	return &AnimalController{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
		transferService:     transferService,
	}
}

// AnimalResponse represents the response for an animal
type AnimalResponse struct {
	ID           string    `json:"id"`
	Species      string    `json:"species"`
	Name         string    `json:"name"`
	BirthDate    time.Time `json:"birthDate"`
	Gender       string    `json:"gender"`
	FavoriteFood string    `json:"favoriteFood"`
	HealthStatus string    `json:"healthStatus"`
	EnclosureID  string    `json:"enclosureId,omitempty"`
}

// CreateAnimalRequest represents the request to create an animal
type CreateAnimalRequest struct {
	ID           string    `json:"id"`
	Species      string    `json:"species"`
	Name         string    `json:"name"`
	BirthDate    time.Time `json:"birthDate"`
	Gender       string    `json:"gender"`
	FavoriteFood string    `json:"favoriteFood"`
	HealthStatus string    `json:"healthStatus"`
}

// TransferAnimalRequest represents the request to transfer an animal
type TransferAnimalRequest struct {
	EnclosureID string `json:"enclosureId"`
}

// RegisterRoutes registers the routes for the AnimalController
func (c *AnimalController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/animals", c.GetAllAnimals)
	mux.HandleFunc("GET /api/animals/{id}", c.GetAnimal)
	mux.HandleFunc("POST /api/animals", c.CreateAnimal)
	mux.HandleFunc("DELETE /api/animals/{id}", c.DeleteAnimal)
	mux.HandleFunc("POST /api/animals/{id}/transfer", c.TransferAnimal)
	mux.HandleFunc("POST /api/animals/{id}/treat", c.TreatAnimal)
}

// GetAllAnimals handles GET /api/animals
func (c *AnimalController) GetAllAnimals(w http.ResponseWriter, r *http.Request) {
	animals, err := c.animalRepository.GetAll()
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

// GetAnimal handles GET /api/animals/{id}
func (c *AnimalController) GetAnimal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	animal, err := c.animalRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := AnimalResponse{
		ID:           animal.ID,
		Species:      animal.Species.String(),
		Name:         animal.Name.String(),
		BirthDate:    animal.BirthDate.Time(),
		Gender:       string(animal.Gender),
		FavoriteFood: string(animal.FavoriteFood),
		HealthStatus: string(animal.HealthStatus),
		EnclosureID:  animal.EnclosureID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateAnimal handles POST /api/animals
func (c *AnimalController) CreateAnimal(w http.ResponseWriter, r *http.Request) {
	var request CreateAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create value objects
	species, err := domain.NewSpecies(request.Species)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name, err := domain.NewAnimalName(request.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	birthDate, err := domain.NewBirthDate(request.BirthDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create the animal
	animal, err := domain.NewAnimal(
		request.ID,
		species,
		name,
		birthDate,
		domain.Gender(request.Gender),
		domain.FoodType(request.FavoriteFood),
		domain.HealthStatus(request.HealthStatus),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the animal
	if err := c.animalRepository.Save(animal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := AnimalResponse{
		ID:           animal.ID,
		Species:      animal.Species.String(),
		Name:         animal.Name.String(),
		BirthDate:    animal.BirthDate.Time(),
		Gender:       string(animal.Gender),
		FavoriteFood: string(animal.FavoriteFood),
		HealthStatus: string(animal.HealthStatus),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// DeleteAnimal handles DELETE /api/animals/{id}
func (c *AnimalController) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	// Check if the animal exists
	animal, err := c.animalRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// If the animal is in an enclosure, remove it first
	if animal.IsInEnclosure() {
		enclosure, err := c.enclosureRepository.GetByID(animal.EnclosureID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := enclosure.RemoveAnimal(animal.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := c.enclosureRepository.Save(enclosure); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Delete the animal
	if err := c.animalRepository.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TransferAnimal handles POST /api/animals/{id}/transfer
func (c *AnimalController) TransferAnimal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	var request TransferAnimalRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.EnclosureID == "" {
		http.Error(w, "Enclosure ID is required", http.StatusBadRequest)
		return
	}

	// Transfer the animal
	if err := c.transferService.TransferAnimal(id, request.EnclosureID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// TreatAnimal handles POST /api/animals/{id}/treat
func (c *AnimalController) TreatAnimal(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Animal ID is required", http.StatusBadRequest)
		return
	}

	// Get the animal
	animal, err := c.animalRepository.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Treat the animal
	if err := animal.Treat(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the animal
	if err := c.animalRepository.Save(animal); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}