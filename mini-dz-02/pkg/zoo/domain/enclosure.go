package domain

import (
	"errors"
	"fmt"
)

// Enclosure represents an animal enclosure in the zoo
type Enclosure struct {
	ID               string
	Type             EnclosureType
	Size             EnclosureSize
	MaxCapacity      Capacity
	CurrentAnimalIDs []string
	CleaningStatus   bool // true if clean, false if needs cleaning
}

// NewEnclosure creates a new Enclosure with validation
func NewEnclosure(
	id string,
	enclosureType EnclosureType,
	size EnclosureSize,
	maxCapacity Capacity,
) (*Enclosure, error) {
	if id == "" {
		return nil, errors.New("enclosure ID cannot be empty")
	}

	if size.Value() < maxCapacity.Value() {
		return nil, errors.New("enclosure size must be greater than or equal to max capacity")
	}

	return &Enclosure{
		ID:               id,
		Type:             enclosureType,
		Size:             size,
		MaxCapacity:      maxCapacity,
		CurrentAnimalIDs: make([]string, 0),
		CleaningStatus:   true, // New enclosures are clean
	}, nil
}

// AddAnimal adds an animal to the enclosure
func (e *Enclosure) AddAnimal(animalID string) error {
	if animalID == "" {
		return errors.New("animal ID cannot be empty")
	}

	// Check if the animal is already in this enclosure
	for _, id := range e.CurrentAnimalIDs {
		if id == animalID {
			return errors.New("animal is already in this enclosure")
		}
	}

	// Check if the enclosure is at capacity
	if len(e.CurrentAnimalIDs) >= e.MaxCapacity.Value() {
		return errors.New("enclosure is at maximum capacity")
	}

	e.CurrentAnimalIDs = append(e.CurrentAnimalIDs, animalID)
	e.CleaningStatus = false // Adding an animal makes the enclosure dirty

	fmt.Printf("Animal %s has been added to enclosure %s\n", animalID, e.ID)
	return nil
}

// RemoveAnimal removes an animal from the enclosure
func (e *Enclosure) RemoveAnimal(animalID string) error {
	if animalID == "" {
		return errors.New("animal ID cannot be empty")
	}

	// Find the animal in the enclosure
	index := -1
	for i, id := range e.CurrentAnimalIDs {
		if id == animalID {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("animal is not in this enclosure")
	}

	// Remove the animal from the slice
	e.CurrentAnimalIDs = append(e.CurrentAnimalIDs[:index], e.CurrentAnimalIDs[index+1:]...)

	fmt.Printf("Animal %s has been removed from enclosure %s\n", animalID, e.ID)
	return nil
}

// Clean cleans the enclosure
func (e *Enclosure) Clean() error {
	if e.CleaningStatus {
		return errors.New("enclosure is already clean")
	}

	e.CleaningStatus = true
	fmt.Printf("Enclosure %s has been cleaned\n", e.ID)
	return nil
}

// IsClean returns whether the enclosure is clean
func (e *Enclosure) IsClean() bool {
	return e.CleaningStatus
}

// CurrentAnimalCount returns the current number of animals in the enclosure
func (e *Enclosure) CurrentAnimalCount() int {
	return len(e.CurrentAnimalIDs)
}

// HasSpace returns whether the enclosure has space for more animals
func (e *Enclosure) HasSpace() bool {
	return len(e.CurrentAnimalIDs) < e.MaxCapacity.Value()
}

// IsEmpty returns whether the enclosure is empty
func (e *Enclosure) IsEmpty() bool {
	return len(e.CurrentAnimalIDs) == 0
}