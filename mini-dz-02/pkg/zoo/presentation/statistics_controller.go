package presentation

import (
	"encoding/json"
	"net/http"

	"mini-dz-02/pkg/zoo/application"
)

// StatisticsController handles HTTP requests related to zoo statistics
type StatisticsController struct {
	statisticsService *application.ZooStatisticsService
}

// NewStatisticsController creates a new StatisticsController
func NewStatisticsController(
	statisticsService *application.ZooStatisticsService,
) *StatisticsController {
	return &StatisticsController{
		statisticsService: statisticsService,
	}
}

// ZooStatisticsResponse represents the response for zoo statistics
type ZooStatisticsResponse struct {
	TotalAnimals        int               `json:"totalAnimals"`
	HealthyAnimals      int               `json:"healthyAnimals"`
	SickAnimals         int               `json:"sickAnimals"`
	TotalEnclosures     int               `json:"totalEnclosures"`
	AvailableEnclosures int               `json:"availableEnclosures"`
	FullEnclosures      int               `json:"fullEnclosures"`
	EmptyEnclosures     int               `json:"emptyEnclosures"`
	AnimalsBySpecies    map[string]int    `json:"animalsBySpecies"`
	AnimalsByGender     map[string]int    `json:"animalsByGender"`
	AnimalsByEnclosure  map[string]int    `json:"animalsByEnclosure"`
}

// RegisterRoutes registers the routes for the StatisticsController
func (c *StatisticsController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/statistics", c.GetZooStatistics)
	mux.HandleFunc("GET /api/statistics/species", c.GetAnimalCountBySpecies)
	mux.HandleFunc("GET /api/statistics/enclosure-utilization", c.GetEnclosureUtilization)
	mux.HandleFunc("GET /api/statistics/health", c.GetHealthStatusStatistics)
}

// GetZooStatistics handles GET /api/statistics
func (c *StatisticsController) GetZooStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := c.statisticsService.GetZooStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert gender map keys from domain.Gender to string
	animalsByGender := make(map[string]int)
	for gender, count := range stats.AnimalsByGender {
		animalsByGender[string(gender)] = count
	}

	response := ZooStatisticsResponse{
		TotalAnimals:        stats.TotalAnimals,
		HealthyAnimals:      stats.HealthyAnimals,
		SickAnimals:         stats.SickAnimals,
		TotalEnclosures:     stats.TotalEnclosures,
		AvailableEnclosures: stats.AvailableEnclosures,
		FullEnclosures:      stats.FullEnclosures,
		EmptyEnclosures:     stats.EmptyEnclosures,
		AnimalsBySpecies:    stats.AnimalsBySpecies,
		AnimalsByGender:     animalsByGender,
		AnimalsByEnclosure:  stats.AnimalsByEnclosure,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAnimalCountBySpecies handles GET /api/statistics/species
func (c *StatisticsController) GetAnimalCountBySpecies(w http.ResponseWriter, r *http.Request) {
	animalsBySpecies, err := c.statisticsService.GetAnimalCountBySpecies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animalsBySpecies)
}

// GetEnclosureUtilization handles GET /api/statistics/enclosure-utilization
func (c *StatisticsController) GetEnclosureUtilization(w http.ResponseWriter, r *http.Request) {
	utilization, err := c.statisticsService.GetEnclosureUtilization()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(utilization)
}

// GetHealthStatusStatistics handles GET /api/statistics/health
func (c *StatisticsController) GetHealthStatusStatistics(w http.ResponseWriter, r *http.Request) {
	healthStats, err := c.statisticsService.GetHealthStatusStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert health status map keys from domain.HealthStatus to string
	healthStatsResponse := make(map[string]int)
	for status, count := range healthStats {
		healthStatsResponse[string(status)] = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthStatsResponse)
}