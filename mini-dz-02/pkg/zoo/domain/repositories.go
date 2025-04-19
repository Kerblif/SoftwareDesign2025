package domain

// AnimalRepository defines the interface for animal data access
type AnimalRepository interface {
	GetByID(id string) (*Animal, error)
	GetAll() ([]*Animal, error)
	Save(animal *Animal) error
	Delete(id string) error
	GetByEnclosureID(enclosureID string) ([]*Animal, error)
}

// EnclosureRepository defines the interface for enclosure data access
type EnclosureRepository interface {
	GetByID(id string) (*Enclosure, error)
	GetAll() ([]*Enclosure, error)
	Save(enclosure *Enclosure) error
	Delete(id string) error
	GetByType(enclosureType EnclosureType) ([]*Enclosure, error)
	GetAvailable() ([]*Enclosure, error)
}

// FeedingScheduleRepository defines the interface for feeding schedule data access
type FeedingScheduleRepository interface {
	GetByID(id string) (*FeedingSchedule, error)
	GetAll() ([]*FeedingSchedule, error)
	Save(schedule *FeedingSchedule) error
	Delete(id string) error
	GetByAnimalID(animalID string) ([]*FeedingSchedule, error)
	GetDueSchedules() ([]*FeedingSchedule, error)
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(event Event) error
}