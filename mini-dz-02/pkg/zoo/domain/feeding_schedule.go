package domain

import (
	"errors"
	"fmt"
	"time"
)

// FeedingSchedule represents a feeding schedule for an animal
type FeedingSchedule struct {
	ID          string
	AnimalID    string
	FeedingTime FeedingTime
	FoodType    FoodType
	Completed   bool
}

// NewFeedingSchedule creates a new FeedingSchedule with validation
func NewFeedingSchedule(
	id string,
	animalID string,
	feedingTime FeedingTime,
	foodType FoodType,
) (*FeedingSchedule, error) {
	if id == "" {
		return nil, errors.New("feeding schedule ID cannot be empty")
	}

	if animalID == "" {
		return nil, errors.New("animal ID cannot be empty")
	}

	return &FeedingSchedule{
		ID:          id,
		AnimalID:    animalID,
		FeedingTime: feedingTime,
		FoodType:    foodType,
		Completed:   false,
	}, nil
}

// ChangeSchedule changes the feeding time and food type
func (fs *FeedingSchedule) ChangeSchedule(newFeedingTime FeedingTime, newFoodType FoodType) error {
	if fs.Completed {
		return errors.New("cannot change a completed feeding schedule")
	}

	fs.FeedingTime = newFeedingTime
	fs.FoodType = newFoodType

	fmt.Printf("Feeding schedule for animal %s has been updated\n", fs.AnimalID)
	return nil
}

// MarkCompleted marks the feeding schedule as completed
func (fs *FeedingSchedule) MarkCompleted() error {
	if fs.Completed {
		return errors.New("feeding schedule is already marked as completed")
	}

	fs.Completed = true
	fmt.Printf("Feeding schedule for animal %s has been marked as completed\n", fs.AnimalID)
	return nil
}

// IsCompleted returns whether the feeding schedule is completed
func (fs *FeedingSchedule) IsCompleted() bool {
	return fs.Completed
}

// IsDue checks if the feeding is due (current time is after feeding time)
func (fs *FeedingSchedule) IsDue() bool {
	return time.Now().After(fs.FeedingTime.Time())
}

// TimeUntilFeeding returns the duration until the feeding time
func (fs *FeedingSchedule) TimeUntilFeeding() time.Duration {
	now := time.Now()
	feedingTime := fs.FeedingTime.Time()
	
	if now.After(feedingTime) {
		return 0
	}
	
	return feedingTime.Sub(now)
}