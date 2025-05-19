package domain

import (
	"errors"
	"time"
)

// Gender represents the gender of an animal
type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

// HealthStatus represents the health status of an animal
type HealthStatus string

const (
	Healthy HealthStatus = "healthy"
	Sick    HealthStatus = "sick"
)

// EnclosureType represents the type of enclosure
type EnclosureType string

const (
	Predator   EnclosureType = "predator"
	Herbivore  EnclosureType = "herbivore"
	Aviary     EnclosureType = "aviary"
	Aquarium   EnclosureType = "aquarium"
	Terrarium  EnclosureType = "terrarium"
)

// FoodType represents the type of food
type FoodType string

const (
	Meat      FoodType = "meat"
	Vegetables FoodType = "vegetables"
	Fruits    FoodType = "fruits"
	Insects   FoodType = "insects"
	Seeds     FoodType = "seeds"
)

// Species represents an animal species with validation
type Species struct {
	value string
}

// NewSpecies creates a new Species value object with validation
func NewSpecies(value string) (Species, error) {
	if value == "" {
		return Species{}, errors.New("species cannot be empty")
	}
	return Species{value: value}, nil
}

// String returns the string representation of Species
func (s Species) String() string {
	return s.value
}

// AnimalName represents an animal's name with validation
type AnimalName struct {
	value string
}

// NewAnimalName creates a new AnimalName value object with validation
func NewAnimalName(value string) (AnimalName, error) {
	if value == "" {
		return AnimalName{}, errors.New("animal name cannot be empty")
	}
	return AnimalName{value: value}, nil
}

// String returns the string representation of AnimalName
func (n AnimalName) String() string {
	return n.value
}

// BirthDate represents an animal's birth date with validation
type BirthDate struct {
	value time.Time
}

// NewBirthDate creates a new BirthDate value object with validation
func NewBirthDate(value time.Time) (BirthDate, error) {
	if value.After(time.Now()) {
		return BirthDate{}, errors.New("birth date cannot be in the future")
	}
	return BirthDate{value: value}, nil
}

// Time returns the time.Time representation of BirthDate
func (d BirthDate) Time() time.Time {
	return d.value
}

// EnclosureSize represents the size of an enclosure with validation
type EnclosureSize struct {
	value int
}

// NewEnclosureSize creates a new EnclosureSize value object with validation
func NewEnclosureSize(value int) (EnclosureSize, error) {
	if value <= 0 {
		return EnclosureSize{}, errors.New("enclosure size must be positive")
	}
	return EnclosureSize{value: value}, nil
}

// Value returns the int representation of EnclosureSize
func (s EnclosureSize) Value() int {
	return s.value
}

// Capacity represents the capacity of an enclosure with validation
type Capacity struct {
	value int
}

// NewCapacity creates a new Capacity value object with validation
func NewCapacity(value int) (Capacity, error) {
	if value <= 0 {
		return Capacity{}, errors.New("capacity must be positive")
	}
	return Capacity{value: value}, nil
}

// Value returns the int representation of Capacity
func (c Capacity) Value() int {
	return c.value
}

// FeedingTime represents a feeding time with validation
type FeedingTime struct {
	value time.Time
}

// NewFeedingTime creates a new FeedingTime value object with validation
func NewFeedingTime(value time.Time) (FeedingTime, error) {
	return FeedingTime{value: value}, nil
}

// Time returns the time.Time representation of FeedingTime
func (t FeedingTime) Time() time.Time {
	return t.value
}