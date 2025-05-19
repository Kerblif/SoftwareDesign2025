package domain

import (
	"errors"
	"fmt"
)

// Animal represents an animal in the zoo
type Animal struct {
	ID          string
	Species     Species
	Name        AnimalName
	BirthDate   BirthDate
	Gender      Gender
	FavoriteFood FoodType
	HealthStatus HealthStatus
	EnclosureID  string
}

// NewAnimal creates a new Animal with validation
func NewAnimal(
	id string,
	species Species,
	name AnimalName,
	birthDate BirthDate,
	gender Gender,
	favoriteFood FoodType,
	healthStatus HealthStatus,
) (*Animal, error) {
	if id == "" {
		return nil, errors.New("animal ID cannot be empty")
	}

	return &Animal{
		ID:           id,
		Species:      species,
		Name:         name,
		BirthDate:    birthDate,
		Gender:       gender,
		FavoriteFood: favoriteFood,
		HealthStatus: healthStatus,
	}, nil
}

// Feed feeds the animal with the specified food type
func (a *Animal) Feed(foodType FoodType) error {
	if a.HealthStatus == Sick {
		return errors.New("cannot feed a sick animal, treat it first")
	}
	
	// In a real application, we might have more complex logic here
	// For example, checking if the food is appropriate for this animal species
	
	fmt.Printf("Animal %s has been fed with %s\n", a.Name.String(), foodType)
	return nil
}

// Treat treats a sick animal
func (a *Animal) Treat() error {
	if a.HealthStatus == Healthy {
		return errors.New("animal is already healthy")
	}
	
	a.HealthStatus = Healthy
	fmt.Printf("Animal %s has been treated and is now healthy\n", a.Name.String())
	return nil
}

// MoveToEnclosure moves the animal to a different enclosure
func (a *Animal) MoveToEnclosure(enclosureID string) error {
	if enclosureID == "" {
		return errors.New("enclosure ID cannot be empty")
	}
	
	if a.EnclosureID == enclosureID {
		return errors.New("animal is already in this enclosure")
	}
	
	oldEnclosureID := a.EnclosureID
	a.EnclosureID = enclosureID
	
	// In a real application, we would publish the AnimalMovedEvent here
	// For simplicity, we'll just print a message
	fmt.Printf("Animal %s has been moved from enclosure %s to enclosure %s\n", 
		a.Name.String(), oldEnclosureID, enclosureID)
	
	return nil
}

// SetHealthStatus updates the health status of the animal
func (a *Animal) SetHealthStatus(status HealthStatus) {
	a.HealthStatus = status
}

// IsInEnclosure checks if the animal is in an enclosure
func (a *Animal) IsInEnclosure() bool {
	return a.EnclosureID != ""
}