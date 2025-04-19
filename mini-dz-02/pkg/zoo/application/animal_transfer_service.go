package application

import (
	"errors"
	"fmt"

	"mini-dz-02/pkg/zoo/domain"
)

// AnimalTransferService handles the logic for transferring animals between enclosures
type AnimalTransferService struct {
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
	eventPublisher      domain.EventPublisher
}

// NewAnimalTransferService creates a new AnimalTransferService
func NewAnimalTransferService(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
	eventPublisher domain.EventPublisher,
) *AnimalTransferService {
	return &AnimalTransferService{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
		eventPublisher:      eventPublisher,
	}
}

// TransferAnimal transfers an animal from one enclosure to another
func (s *AnimalTransferService) TransferAnimal(animalID, targetEnclosureID string) error {
	// Get the animal
	animal, err := s.animalRepository.GetByID(animalID)
	if err != nil {
		return fmt.Errorf("failed to get animal: %w", err)
	}

	// Get the target enclosure
	targetEnclosure, err := s.enclosureRepository.GetByID(targetEnclosureID)
	if err != nil {
		return fmt.Errorf("failed to get target enclosure: %w", err)
	}

	// Check if the target enclosure has space
	if !targetEnclosure.HasSpace() {
		return errors.New("target enclosure is at maximum capacity")
	}

	// Check if the animal is already in an enclosure
	fromEnclosureID := animal.EnclosureID
	if fromEnclosureID != "" {
		// Get the current enclosure
		fromEnclosure, err := s.enclosureRepository.GetByID(fromEnclosureID)
		if err != nil {
			return fmt.Errorf("failed to get current enclosure: %w", err)
		}

		// Remove the animal from the current enclosure
		if err := fromEnclosure.RemoveAnimal(animalID); err != nil {
			return fmt.Errorf("failed to remove animal from current enclosure: %w", err)
		}

		// Save the updated enclosure
		if err := s.enclosureRepository.Save(fromEnclosure); err != nil {
			return fmt.Errorf("failed to save current enclosure: %w", err)
		}
	}

	// Add the animal to the target enclosure
	if err := targetEnclosure.AddAnimal(animalID); err != nil {
		return fmt.Errorf("failed to add animal to target enclosure: %w", err)
	}

	// Update the animal's enclosure ID
	if err := animal.MoveToEnclosure(targetEnclosureID); err != nil {
		return fmt.Errorf("failed to update animal's enclosure ID: %w", err)
	}

	// Save the updated animal and enclosure
	if err := s.animalRepository.Save(animal); err != nil {
		return fmt.Errorf("failed to save animal: %w", err)
	}

	if err := s.enclosureRepository.Save(targetEnclosure); err != nil {
		return fmt.Errorf("failed to save target enclosure: %w", err)
	}

	// Publish the AnimalMovedEvent
	event := domain.NewAnimalMovedEvent(animalID, fromEnclosureID, targetEnclosureID)
	if err := s.eventPublisher.Publish(event); err != nil {
		return fmt.Errorf("failed to publish AnimalMovedEvent: %w", err)
	}

	return nil
}

// GetAvailableEnclosures returns a list of enclosures that have space for more animals
func (s *AnimalTransferService) GetAvailableEnclosures() ([]*domain.Enclosure, error) {
	return s.enclosureRepository.GetAvailable()
}

// GetAnimalsByEnclosure returns a list of animals in a specific enclosure
func (s *AnimalTransferService) GetAnimalsByEnclosure(enclosureID string) ([]*domain.Animal, error) {
	return s.animalRepository.GetByEnclosureID(enclosureID)
}

// GetAnimalByID returns an animal by its ID
func (s *AnimalTransferService) GetAnimalByID(animalID string) (*domain.Animal, error) {
	animal, err := s.animalRepository.GetByID(animalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get animal: %w", err)
	}
	return animal, nil
}

// GetAllAnimals returns all animals
func (s *AnimalTransferService) GetAllAnimals() ([]*domain.Animal, error) {
	animals, err := s.animalRepository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all animals: %w", err)
	}
	return animals, nil
}

// CreateAnimal creates a new animal
func (s *AnimalTransferService) CreateAnimal(
	id string,
	species domain.Species,
	name domain.AnimalName,
	birthDate domain.BirthDate,
	gender domain.Gender,
	favoriteFood domain.FoodType,
	healthStatus domain.HealthStatus,
) (*domain.Animal, error) {
	// Create the animal
	animal, err := domain.NewAnimal(
		id,
		species,
		name,
		birthDate,
		gender,
		favoriteFood,
		healthStatus,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create animal: %w", err)
	}

	// Save the animal
	if err := s.animalRepository.Save(animal); err != nil {
		return nil, fmt.Errorf("failed to save animal: %w", err)
	}

	return animal, nil
}

// DeleteAnimal deletes an animal
func (s *AnimalTransferService) DeleteAnimal(animalID string) error {
	// Check if the animal exists
	animal, err := s.animalRepository.GetByID(animalID)
	if err != nil {
		return fmt.Errorf("failed to get animal: %w", err)
	}

	// If the animal is in an enclosure, remove it first
	if animal.IsInEnclosure() {
		enclosure, err := s.enclosureRepository.GetByID(animal.EnclosureID)
		if err != nil {
			return fmt.Errorf("failed to get enclosure: %w", err)
		}

		if err := enclosure.RemoveAnimal(animal.ID); err != nil {
			return fmt.Errorf("failed to remove animal from enclosure: %w", err)
		}

		if err := s.enclosureRepository.Save(enclosure); err != nil {
			return fmt.Errorf("failed to save enclosure: %w", err)
		}
	}

	// Delete the animal
	if err := s.animalRepository.Delete(animalID); err != nil {
		return fmt.Errorf("failed to delete animal: %w", err)
	}

	return nil
}

// TreatAnimal treats a sick animal
func (s *AnimalTransferService) TreatAnimal(animalID string) error {
	// Get the animal
	animal, err := s.animalRepository.GetByID(animalID)
	if err != nil {
		return fmt.Errorf("failed to get animal: %w", err)
	}

	// Treat the animal
	if err := animal.Treat(); err != nil {
		return fmt.Errorf("failed to treat animal: %w", err)
	}

	// Save the animal
	if err := s.animalRepository.Save(animal); err != nil {
		return fmt.Errorf("failed to save animal: %w", err)
	}

	return nil
}
