package domain_test

import (
	"testing"
	"time"

	"mini-dz-02/pkg/zoo/domain"
)

// Helper function to create a valid feeding schedule for testing
func createValidFeedingSchedule(t *testing.T) *domain.FeedingSchedule {
	// Create a feeding time 1 hour in the future
	feedingTime, err := domain.NewFeedingTime(time.Now().Add(1 * time.Hour))
	if err != nil {
		t.Fatalf("Failed to create feeding time: %v", err)
	}

	schedule, err := domain.NewFeedingSchedule(
		"schedule-1",
		"animal-1",
		feedingTime,
		domain.Meat,
	)
	if err != nil {
		t.Fatalf("Failed to create feeding schedule: %v", err)
	}

	return schedule
}

// Helper function to create a feeding time in the past
func createPastFeedingTime(t *testing.T) domain.FeedingTime {
	feedingTime, err := domain.NewFeedingTime(time.Now().Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("Failed to create past feeding time: %v", err)
	}
	return feedingTime
}

// Helper function to create a feeding time in the future
func createFutureFeedingTime(t *testing.T) domain.FeedingTime {
	feedingTime, err := domain.NewFeedingTime(time.Now().Add(1 * time.Hour))
	if err != nil {
		t.Fatalf("Failed to create future feeding time: %v", err)
	}
	return feedingTime
}

func TestNewFeedingSchedule(t *testing.T) {
	// Test valid feeding schedule creation
	feedingTime := createFutureFeedingTime(t)

	schedule, err := domain.NewFeedingSchedule(
		"schedule-1",
		"animal-1",
		feedingTime,
		domain.Meat,
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if schedule.ID != "schedule-1" {
		t.Errorf("Expected ID to be 'schedule-1', got %s", schedule.ID)
	}

	if schedule.AnimalID != "animal-1" {
		t.Errorf("Expected AnimalID to be 'animal-1', got %s", schedule.AnimalID)
	}

	if schedule.FoodType != domain.Meat {
		t.Errorf("Expected FoodType to be Meat, got %s", schedule.FoodType)
	}

	if schedule.Completed {
		t.Errorf("Expected Completed to be false, got true")
	}

	// Test invalid feeding schedule creation (empty ID)
	_, err = domain.NewFeedingSchedule(
		"",
		"animal-1",
		feedingTime,
		domain.Meat,
	)

	if err == nil {
		t.Error("Expected error for empty ID, got nil")
	}

	// Test invalid feeding schedule creation (empty animal ID)
	_, err = domain.NewFeedingSchedule(
		"schedule-1",
		"",
		feedingTime,
		domain.Meat,
	)

	if err == nil {
		t.Error("Expected error for empty animal ID, got nil")
	}
}

func TestFeedingSchedule_ChangeSchedule(t *testing.T) {
	schedule := createValidFeedingSchedule(t)
	originalFeedingTime := schedule.FeedingTime
	originalFoodType := schedule.FoodType

	// Test changing the schedule
	newFeedingTime := createFutureFeedingTime(t)
	newFoodType := domain.Vegetables

	err := schedule.ChangeSchedule(newFeedingTime, newFoodType)
	if err != nil {
		t.Errorf("Expected no error when changing schedule, got %v", err)
	}

	if schedule.FeedingTime == originalFeedingTime {
		t.Errorf("Expected feeding time to change, but it remained the same")
	}

	if schedule.FoodType == originalFoodType {
		t.Errorf("Expected food type to change, but it remained the same")
	}

	if schedule.FoodType != newFoodType {
		t.Errorf("Expected food type to be %s, got %s", newFoodType, schedule.FoodType)
	}

	// Test changing a completed schedule
	schedule.Completed = true
	err = schedule.ChangeSchedule(createFutureFeedingTime(t), domain.Fruits)
	if err == nil {
		t.Error("Expected error when changing a completed schedule, got nil")
	}
}

func TestFeedingSchedule_MarkCompleted(t *testing.T) {
	schedule := createValidFeedingSchedule(t)

	// Test marking as completed
	err := schedule.MarkCompleted()
	if err != nil {
		t.Errorf("Expected no error when marking as completed, got %v", err)
	}

	if !schedule.Completed {
		t.Errorf("Expected Completed to be true after marking as completed, got false")
	}

	// Test marking an already completed schedule
	err = schedule.MarkCompleted()
	if err == nil {
		t.Error("Expected error when marking an already completed schedule, got nil")
	}
}

func TestFeedingSchedule_IsCompleted(t *testing.T) {
	schedule := createValidFeedingSchedule(t)

	// Test initial state (not completed)
	if schedule.IsCompleted() {
		t.Error("Expected new schedule to not be completed")
	}

	// Test after marking as completed
	schedule.Completed = true
	if !schedule.IsCompleted() {
		t.Error("Expected schedule to be completed after setting Completed to true")
	}
}

func TestFeedingSchedule_IsDue(t *testing.T) {
	// Test a schedule that is not due yet (future feeding time)
	futureFeedingTime := createFutureFeedingTime(t)
	futureSchedule, _ := domain.NewFeedingSchedule(
		"schedule-1",
		"animal-1",
		futureFeedingTime,
		domain.Meat,
	)

	if futureSchedule.IsDue() {
		t.Error("Expected future schedule to not be due yet")
	}

	// Test a schedule that is due (past feeding time)
	pastFeedingTime := createPastFeedingTime(t)
	pastSchedule, _ := domain.NewFeedingSchedule(
		"schedule-2",
		"animal-1",
		pastFeedingTime,
		domain.Meat,
	)

	if !pastSchedule.IsDue() {
		t.Error("Expected past schedule to be due")
	}
}

func TestFeedingSchedule_TimeUntilFeeding(t *testing.T) {
	// Test a schedule with future feeding time
	futureFeedingTime, _ := domain.NewFeedingTime(time.Now().Add(1 * time.Hour))
	futureSchedule, _ := domain.NewFeedingSchedule(
		"schedule-1",
		"animal-1",
		futureFeedingTime,
		domain.Meat,
	)

	timeUntil := futureSchedule.TimeUntilFeeding()
	if timeUntil <= 0 {
		t.Errorf("Expected positive duration for future schedule, got %v", timeUntil)
	}

	// The exact duration will vary, but it should be less than or equal to 1 hour
	if timeUntil > 1*time.Hour {
		t.Errorf("Expected duration to be less than 1 hour, got %v", timeUntil)
	}

	// Test a schedule with past feeding time
	pastFeedingTime, _ := domain.NewFeedingTime(time.Now().Add(-1 * time.Hour))
	pastSchedule, _ := domain.NewFeedingSchedule(
		"schedule-2",
		"animal-1",
		pastFeedingTime,
		domain.Meat,
	)

	timeUntil = pastSchedule.TimeUntilFeeding()
	if timeUntil != 0 {
		t.Errorf("Expected zero duration for past schedule, got %v", timeUntil)
	}
}
