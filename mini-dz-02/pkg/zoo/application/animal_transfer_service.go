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