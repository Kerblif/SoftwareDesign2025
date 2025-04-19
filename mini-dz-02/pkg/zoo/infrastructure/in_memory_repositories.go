package infrastructure

import (
	"errors"
	"sync"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

// InMemoryAnimalRepository is an in-memory implementation of the AnimalRepository interface
type InMemoryAnimalRepository struct {
	animals map[string]*domain.Animal
	mutex   sync.RWMutex
}

// NewInMemoryAnimalRepository creates a new InMemoryAnimalRepository
func NewInMemoryAnimalRepository() *InMemoryAnimalRepository {
	return &InMemoryAnimalRepository{
		animals: make(map[string]*domain.Animal),
	}
}

// GetByID returns an animal by ID
func (r *InMemoryAnimalRepository) GetByID(id string) (*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	animal, exists := r.animals[id]
	if !exists {
		return nil, errors.New("animal not found")
	}

	return animal, nil
}

// GetAll returns all animals
func (r *InMemoryAnimalRepository) GetAll() ([]*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	animals := make([]*domain.Animal, 0, len(r.animals))
	for _, animal := range r.animals {
		animals = append(animals, animal)
	}

	return animals, nil
}

// Save saves an animal
func (r *InMemoryAnimalRepository) Save(animal *domain.Animal) error {
	if animal == nil {
		return errors.New("animal cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.animals[animal.ID] = animal
	return nil
}

// Delete deletes an animal
func (r *InMemoryAnimalRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.animals[id]; !exists {
		return errors.New("animal not found")
	}

	delete(r.animals, id)
	return nil
}

// GetByEnclosureID returns animals by enclosure ID
func (r *InMemoryAnimalRepository) GetByEnclosureID(enclosureID string) ([]*domain.Animal, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var animals []*domain.Animal
	for _, animal := range r.animals {
		if animal.EnclosureID == enclosureID {
			animals = append(animals, animal)
		}
	}

	return animals, nil
}

// InMemoryEnclosureRepository is an in-memory implementation of the EnclosureRepository interface
type InMemoryEnclosureRepository struct {
	enclosures map[string]*domain.Enclosure
	mutex      sync.RWMutex
}

// NewInMemoryEnclosureRepository creates a new InMemoryEnclosureRepository
func NewInMemoryEnclosureRepository() *InMemoryEnclosureRepository {
	return &InMemoryEnclosureRepository{
		enclosures: make(map[string]*domain.Enclosure),
	}
}

// GetByID returns an enclosure by ID
func (r *InMemoryEnclosureRepository) GetByID(id string) (*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enclosure, exists := r.enclosures[id]
	if !exists {
		return nil, errors.New("enclosure not found")
	}

	return enclosure, nil
}

// GetAll returns all enclosures
func (r *InMemoryEnclosureRepository) GetAll() ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	enclosures := make([]*domain.Enclosure, 0, len(r.enclosures))
	for _, enclosure := range r.enclosures {
		enclosures = append(enclosures, enclosure)
	}

	return enclosures, nil
}

// Save saves an enclosure
func (r *InMemoryEnclosureRepository) Save(enclosure *domain.Enclosure) error {
	if enclosure == nil {
		return errors.New("enclosure cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.enclosures[enclosure.ID] = enclosure
	return nil
}

// Delete deletes an enclosure
func (r *InMemoryEnclosureRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.enclosures[id]; !exists {
		return errors.New("enclosure not found")
	}

	delete(r.enclosures, id)
	return nil
}

// GetByType returns enclosures by type
func (r *InMemoryEnclosureRepository) GetByType(enclosureType domain.EnclosureType) ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var enclosures []*domain.Enclosure
	for _, enclosure := range r.enclosures {
		if enclosure.Type == enclosureType {
			enclosures = append(enclosures, enclosure)
		}
	}

	return enclosures, nil
}

// GetAvailable returns enclosures that have space for more animals
func (r *InMemoryEnclosureRepository) GetAvailable() ([]*domain.Enclosure, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var enclosures []*domain.Enclosure
	for _, enclosure := range r.enclosures {
		if enclosure.HasSpace() {
			enclosures = append(enclosures, enclosure)
		}
	}

	return enclosures, nil
}

// InMemoryFeedingScheduleRepository is an in-memory implementation of the FeedingScheduleRepository interface
type InMemoryFeedingScheduleRepository struct {
	schedules map[string]*domain.FeedingSchedule
	mutex     sync.RWMutex
}

// NewInMemoryFeedingScheduleRepository creates a new InMemoryFeedingScheduleRepository
func NewInMemoryFeedingScheduleRepository() *InMemoryFeedingScheduleRepository {
	return &InMemoryFeedingScheduleRepository{
		schedules: make(map[string]*domain.FeedingSchedule),
	}
}

// GetByID returns a feeding schedule by ID
func (r *InMemoryFeedingScheduleRepository) GetByID(id string) (*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	schedule, exists := r.schedules[id]
	if !exists {
		return nil, errors.New("feeding schedule not found")
	}

	return schedule, nil
}

// GetAll returns all feeding schedules
func (r *InMemoryFeedingScheduleRepository) GetAll() ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	schedules := make([]*domain.FeedingSchedule, 0, len(r.schedules))
	for _, schedule := range r.schedules {
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

// Save saves a feeding schedule
func (r *InMemoryFeedingScheduleRepository) Save(schedule *domain.FeedingSchedule) error {
	if schedule == nil {
		return errors.New("feeding schedule cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.schedules[schedule.ID] = schedule
	return nil
}

// Delete deletes a feeding schedule
func (r *InMemoryFeedingScheduleRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.schedules[id]; !exists {
		return errors.New("feeding schedule not found")
	}

	delete(r.schedules, id)
	return nil
}

// GetByAnimalID returns feeding schedules by animal ID
func (r *InMemoryFeedingScheduleRepository) GetByAnimalID(animalID string) ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	for _, schedule := range r.schedules {
		if schedule.AnimalID == animalID {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// GetDueSchedules returns feeding schedules that are due
func (r *InMemoryFeedingScheduleRepository) GetDueSchedules() ([]*domain.FeedingSchedule, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var schedules []*domain.FeedingSchedule
	now := time.Now()
	for _, schedule := range r.schedules {
		if !schedule.IsCompleted() && schedule.FeedingTime.Time().Before(now) {
			schedules = append(schedules, schedule)
		}
	}

	return schedules, nil
}

// InMemoryEventPublisher is an in-memory implementation of the EventPublisher interface
type InMemoryEventPublisher struct {
	events []domain.Event
	mutex  sync.RWMutex
}

// NewInMemoryEventPublisher creates a new InMemoryEventPublisher
func NewInMemoryEventPublisher() *InMemoryEventPublisher {
	return &InMemoryEventPublisher{
		events: make([]domain.Event, 0),
	}
}

// Publish publishes an event
func (p *InMemoryEventPublisher) Publish(event domain.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.events = append(p.events, event)
	return nil
}

// GetEvents returns all published events
func (p *InMemoryEventPublisher) GetEvents() []domain.Event {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	events := make([]domain.Event, len(p.events))
	copy(events, p.events)
	return events
}