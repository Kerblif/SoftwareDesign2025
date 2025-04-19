package application

import (
	"fmt"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

// FeedingOrganizationService handles the logic for organizing feeding schedules
type FeedingOrganizationService struct {
	animalRepository          domain.AnimalRepository
	feedingScheduleRepository domain.FeedingScheduleRepository
	eventPublisher            domain.EventPublisher
}

// NewFeedingOrganizationService creates a new FeedingOrganizationService
func NewFeedingOrganizationService(
	animalRepository domain.AnimalRepository,
	feedingScheduleRepository domain.FeedingScheduleRepository,
	eventPublisher domain.EventPublisher,
) *FeedingOrganizationService {
	return &FeedingOrganizationService{
		animalRepository:          animalRepository,
		feedingScheduleRepository: feedingScheduleRepository,
		eventPublisher:            eventPublisher,
	}
}

// CreateFeedingSchedule creates a new feeding schedule for an animal
func (s *FeedingOrganizationService) CreateFeedingSchedule(
	scheduleID string,
	animalID string,
	feedingTime time.Time,
	foodType domain.FoodType,
) error {
	// Check if the animal exists
	animal, err := s.animalRepository.GetByID(animalID)
	if err != nil {
		return fmt.Errorf("failed to get animal: %w", err)
	}

	// Create a feeding time value object
	feedingTimeVO, err := domain.NewFeedingTime(feedingTime)
	if err != nil {
		return fmt.Errorf("invalid feeding time: %w", err)
	}

	// Create a new feeding schedule
	schedule, err := domain.NewFeedingSchedule(scheduleID, animalID, feedingTimeVO, foodType)
	if err != nil {
		return fmt.Errorf("failed to create feeding schedule: %w", err)
	}

	// Save the feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return fmt.Errorf("failed to save feeding schedule: %w", err)
	}

	fmt.Printf("Created feeding schedule for animal %s (%s) at %s with food type %s\n",
		animal.Name.String(), animal.Species.String(), feedingTime.Format(time.RFC3339), foodType)

	return nil
}

// UpdateFeedingSchedule updates an existing feeding schedule
func (s *FeedingOrganizationService) UpdateFeedingSchedule(
	scheduleID string,
	newFeedingTime time.Time,
	newFoodType domain.FoodType,
) error {
	// Get the feeding schedule
	schedule, err := s.feedingScheduleRepository.GetByID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get feeding schedule: %w", err)
	}

	// Create a feeding time value object
	feedingTimeVO, err := domain.NewFeedingTime(newFeedingTime)
	if err != nil {
		return fmt.Errorf("invalid feeding time: %w", err)
	}

	// Update the feeding schedule
	if err := schedule.ChangeSchedule(feedingTimeVO, newFoodType); err != nil {
		return fmt.Errorf("failed to update feeding schedule: %w", err)
	}

	// Save the updated feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return fmt.Errorf("failed to save updated feeding schedule: %w", err)
	}

	return nil
}

// MarkFeedingCompleted marks a feeding schedule as completed
func (s *FeedingOrganizationService) MarkFeedingCompleted(scheduleID string) error {
	// Get the feeding schedule
	schedule, err := s.feedingScheduleRepository.GetByID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get feeding schedule: %w", err)
	}

	// Mark the feeding schedule as completed
	if err := schedule.MarkCompleted(); err != nil {
		return fmt.Errorf("failed to mark feeding schedule as completed: %w", err)
	}

	// Save the updated feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return fmt.Errorf("failed to save updated feeding schedule: %w", err)
	}

	// Get the animal
	animal, err := s.animalRepository.GetByID(schedule.AnimalID)
	if err != nil {
		return fmt.Errorf("failed to get animal: %w", err)
	}

	// Feed the animal
	if err := animal.Feed(schedule.FoodType); err != nil {
		return fmt.Errorf("failed to feed animal: %w", err)
	}

	// Save the updated animal
	if err := s.animalRepository.Save(animal); err != nil {
		return fmt.Errorf("failed to save updated animal: %w", err)
	}

	return nil
}

// GetDueFeedings returns a list of feeding schedules that are due
func (s *FeedingOrganizationService) GetDueFeedings() ([]*domain.FeedingSchedule, error) {
	dueSchedules, err := s.feedingScheduleRepository.GetDueSchedules()
	if err != nil {
		return nil, fmt.Errorf("failed to get due feeding schedules: %w", err)
	}

	// Publish FeedingTimeEvent for each due schedule that is not completed
	for _, schedule := range dueSchedules {
		if !schedule.IsCompleted() {
			event := domain.NewFeedingTimeEvent(schedule.AnimalID, schedule.FoodType)
			if err := s.eventPublisher.Publish(event); err != nil {
				fmt.Printf("Warning: failed to publish FeedingTimeEvent for animal %s: %v\n", schedule.AnimalID, err)
			}
		}
	}

	return dueSchedules, nil
}

// GetFeedingSchedulesByAnimal returns a list of feeding schedules for a specific animal
func (s *FeedingOrganizationService) GetFeedingSchedulesByAnimal(animalID string) ([]*domain.FeedingSchedule, error) {
	return s.feedingScheduleRepository.GetByAnimalID(animalID)
}

// GetAllFeedingSchedules returns a list of all feeding schedules
func (s *FeedingOrganizationService) GetAllFeedingSchedules() ([]*domain.FeedingSchedule, error) {
	return s.feedingScheduleRepository.GetAll()
}