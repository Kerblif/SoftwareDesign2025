package application

import (
	"fmt"

	"mini-dz-02/pkg/zoo/domain"
)

// EnclosureManagementService handles the logic for managing enclosures
type EnclosureManagementService struct {
	enclosureRepository domain.EnclosureRepository
	animalRepository    domain.AnimalRepository
	eventPublisher      domain.EventPublisher
}

// NewEnclosureManagementService creates a new EnclosureManagementService
func NewEnclosureManagementService(
	enclosureRepository domain.EnclosureRepository,
	animalRepository domain.AnimalRepository,
	eventPublisher domain.EventPublisher,
) *EnclosureManagementService {
	return &EnclosureManagementService{
		enclosureRepository: enclosureRepository,
		animalRepository:    animalRepository,
		eventPublisher:      eventPublisher,
	}
}

// GetEnclosureByID returns an enclosure by its ID
func (s *EnclosureManagementService) GetEnclosureByID(enclosureID string) (*domain.Enclosure, error) {
	enclosure, err := s.enclosureRepository.GetByID(enclosureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enclosure: %w", err)
	}
	return enclosure, nil
}

// GetAllEnclosures returns all enclosures
func (s *EnclosureManagementService) GetAllEnclosures() ([]*domain.Enclosure, error) {
	enclosures, err := s.enclosureRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all enclosures: %w", err)
	}
	return enclosures, nil
}

// CreateEnclosure creates a new enclosure
func (s *EnclosureManagementService) CreateEnclosure(
	id string,
	enclosureType domain.EnclosureType,
	size domain.EnclosureSize,
	capacity domain.Capacity,
) (*domain.Enclosure, error) {
	// Create the enclosure
	enclosure, err := domain.NewEnclosure(
		id,
		enclosureType,
		size,
		capacity,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create enclosure: %w", err)
	}

	// Save the enclosure
	if err := s.enclosureRepository.Save(enclosure); err != nil {
		return nil, fmt.Errorf("failed to save enclosure: %w", err)
	}

	return enclosure, nil
}

// DeleteEnclosure deletes an enclosure
func (s *EnclosureManagementService) DeleteEnclosure(enclosureID string) error {
	// Check if the enclosure exists
	enclosure, err := s.enclosureRepository.GetByID(enclosureID)
	if err != nil {
		return fmt.Errorf("failed to get enclosure: %w", err)
	}

	// Check if the enclosure is empty
	if len(enclosure.CurrentAnimalIDs) > 0 {
		return fmt.Errorf("cannot delete enclosure with animals")
	}

	// Delete the enclosure
	if err := s.enclosureRepository.Delete(enclosureID); err != nil {
		return fmt.Errorf("failed to delete enclosure: %w", err)
	}

	return nil
}

// CleanEnclosure cleans an enclosure
func (s *EnclosureManagementService) CleanEnclosure(enclosureID string) error {
	// Get the enclosure
	enclosure, err := s.enclosureRepository.GetByID(enclosureID)
	if err != nil {
		return fmt.Errorf("failed to get enclosure: %w", err)
	}

	// Clean the enclosure
	if err := enclosure.Clean(); err != nil {
		return fmt.Errorf("failed to clean enclosure: %w", err)
	}

	// Save the enclosure
	if err := s.enclosureRepository.Save(enclosure); err != nil {
		return fmt.Errorf("failed to save enclosure: %w", err)
	}

	return nil
}

// GetAnimalsInEnclosure returns a list of animals in a specific enclosure
func (s *EnclosureManagementService) GetAnimalsInEnclosure(enclosureID string) ([]*domain.Animal, error) {
	// Get the enclosure
	enclosure, err := s.enclosureRepository.GetByID(enclosureID)
	if err != nil {
		return nil, fmt.Errorf("failed to get enclosure: %w", err)
	}

	// Get the animals in the enclosure
	var animals []*domain.Animal
	for _, animalID := range enclosure.CurrentAnimalIDs {
		animal, err := s.animalRepository.GetByID(animalID)
		if err != nil {
			return nil, fmt.Errorf("failed to get animal: %w", err)
		}
		animals = append(animals, animal)
	}

	return animals, nil
}
