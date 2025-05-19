package application

import (
	"mini-dz-02/pkg/zoo/domain"
)

// ZooStatistics represents statistics about the zoo
type ZooStatistics struct {
	TotalAnimals        int
	HealthyAnimals      int
	SickAnimals         int
	TotalEnclosures     int
	AvailableEnclosures int
	FullEnclosures      int
	EmptyEnclosures     int
	AnimalsBySpecies    map[string]int
	AnimalsByGender     map[domain.Gender]int
	AnimalsByEnclosure  map[string]int
}

// ZooStatisticsService handles the logic for collecting statistics about the zoo
type ZooStatisticsService struct {
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
}

// NewZooStatisticsService creates a new ZooStatisticsService
func NewZooStatisticsService(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
) *ZooStatisticsService {
	return &ZooStatisticsService{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
	}
}

// GetZooStatistics returns statistics about the zoo
func (s *ZooStatisticsService) GetZooStatistics() (*ZooStatistics, error) {
	// Get all animals
	animals, err := s.animalRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// Get all enclosures
	enclosures, err := s.enclosureRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// Initialize statistics
	stats := &ZooStatistics{
		TotalAnimals:        len(animals),
		HealthyAnimals:      0,
		SickAnimals:         0,
		TotalEnclosures:     len(enclosures),
		AvailableEnclosures: 0,
		FullEnclosures:      0,
		EmptyEnclosures:     0,
		AnimalsBySpecies:    make(map[string]int),
		AnimalsByGender:     make(map[domain.Gender]int),
		AnimalsByEnclosure:  make(map[string]int),
	}

	// Count animals by health status, species, gender, and enclosure
	for _, animal := range animals {
		// Count by health status
		if animal.HealthStatus == domain.Healthy {
			stats.HealthyAnimals++
		} else {
			stats.SickAnimals++
		}

		// Count by species
		speciesName := animal.Species.String()
		stats.AnimalsBySpecies[speciesName]++

		// Count by gender
		stats.AnimalsByGender[animal.Gender]++

		// Count by enclosure
		if animal.EnclosureID != "" {
			stats.AnimalsByEnclosure[animal.EnclosureID]++
		}
	}

	// Count enclosures by status
	for _, enclosure := range enclosures {
		if enclosure.IsEmpty() {
			stats.EmptyEnclosures++
		} else if enclosure.HasSpace() {
			stats.AvailableEnclosures++
		} else {
			stats.FullEnclosures++
		}
	}

	return stats, nil
}

// GetAnimalCountBySpecies returns the number of animals of each species
func (s *ZooStatisticsService) GetAnimalCountBySpecies() (map[string]int, error) {
	animals, err := s.animalRepository.GetAll()
	if err != nil {
		return nil, err
	}

	animalsBySpecies := make(map[string]int)
	for _, animal := range animals {
		speciesName := animal.Species.String()
		animalsBySpecies[speciesName]++
	}

	return animalsBySpecies, nil
}

// GetEnclosureUtilization returns the utilization of each enclosure
func (s *ZooStatisticsService) GetEnclosureUtilization() (map[string]float64, error) {
	enclosures, err := s.enclosureRepository.GetAll()
	if err != nil {
		return nil, err
	}

	utilization := make(map[string]float64)
	for _, enclosure := range enclosures {
		if enclosure.MaxCapacity.Value() > 0 {
			utilization[enclosure.ID] = float64(enclosure.CurrentAnimalCount()) / float64(enclosure.MaxCapacity.Value())
		} else {
			utilization[enclosure.ID] = 0
		}
	}

	return utilization, nil
}

// GetHealthStatusStatistics returns statistics about animal health
func (s *ZooStatisticsService) GetHealthStatusStatistics() (map[domain.HealthStatus]int, error) {
	animals, err := s.animalRepository.GetAll()
	if err != nil {
		return nil, err
	}

	healthStats := make(map[domain.HealthStatus]int)
	for _, animal := range animals {
		healthStats[animal.HealthStatus]++
	}

	return healthStats, nil
}