package servers

import (
	"context"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// FeedingScheduleServer implements the FeedingScheduleService
type FeedingScheduleServer struct {
	zoo.UnimplementedFeedingScheduleServiceServer
	feedingScheduleRepository  domain.FeedingScheduleRepository
	animalRepository           domain.AnimalRepository
	feedingOrganizationService *application.FeedingOrganizationService
}

// NewFeedingScheduleServer creates a new FeedingScheduleServer
func NewFeedingScheduleServer(
	feedingScheduleRepository domain.FeedingScheduleRepository,
	animalRepository domain.AnimalRepository,
	feedingOrganizationService *application.FeedingOrganizationService,
) *FeedingScheduleServer {
	return &FeedingScheduleServer{
		feedingScheduleRepository:  feedingScheduleRepository,
		animalRepository:           animalRepository,
		feedingOrganizationService: feedingOrganizationService,
	}
}

// GetFeedingSchedule implements the GetFeedingSchedule method of the FeedingScheduleService
func (s *FeedingScheduleServer) GetFeedingSchedule(ctx context.Context, req *zoo.GetFeedingScheduleRequest) (*zoo.FeedingSchedule, error) {
	schedule, err := s.feedingScheduleRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return ConvertFeedingSchedule(schedule), nil
}

// GetFeedingSchedules implements the GetFeedingSchedules method of the FeedingScheduleService
func (s *FeedingScheduleServer) GetFeedingSchedules(ctx context.Context, req *zoo.Empty) (*zoo.GetFeedingSchedulesResponse, error) {
	schedules, err := s.feedingScheduleRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var protoSchedules []*zoo.FeedingSchedule
	for _, schedule := range schedules {
		protoSchedules = append(protoSchedules, ConvertFeedingSchedule(schedule))
	}

	return &zoo.GetFeedingSchedulesResponse{
		FeedingSchedules: protoSchedules,
	}, nil
}

// CreateFeedingSchedule implements the CreateFeedingSchedule method of the FeedingScheduleService
func (s *FeedingScheduleServer) CreateFeedingSchedule(ctx context.Context, req *zoo.CreateFeedingScheduleRequest) (*zoo.FeedingSchedule, error) {
	// Check if the animal exists
	_, err := s.animalRepository.GetByID(req.AnimalId)
	if err != nil {
		return nil, err
	}

	// Create the feeding time value object
	feedingTime, err := domain.NewFeedingTime(req.FeedingTime.AsTime())
	if err != nil {
		return nil, err
	}

	// Create the feeding schedule
	schedule, err := domain.NewFeedingSchedule(
		req.Id,
		req.AnimalId,
		feedingTime,
		ConvertProtoFoodType(req.FoodType),
	)
	if err != nil {
		return nil, err
	}

	// Save the feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return nil, err
	}

	return ConvertFeedingSchedule(schedule), nil
}

// DeleteFeedingSchedule implements the DeleteFeedingSchedule method of the FeedingScheduleService
func (s *FeedingScheduleServer) DeleteFeedingSchedule(ctx context.Context, req *zoo.DeleteFeedingScheduleRequest) (*zoo.Empty, error) {
	// Check if the feeding schedule exists
	_, err := s.feedingScheduleRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Delete the feeding schedule
	if err := s.feedingScheduleRepository.Delete(req.Id); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// UpdateFeedingSchedule implements the UpdateFeedingSchedule method of the FeedingScheduleService
func (s *FeedingScheduleServer) UpdateFeedingSchedule(ctx context.Context, req *zoo.UpdateFeedingScheduleRequest) (*zoo.FeedingSchedule, error) {
	// Get the feeding schedule
	schedule, err := s.feedingScheduleRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Create the feeding time value object
	feedingTime, err := domain.NewFeedingTime(req.FeedingTime.AsTime())
	if err != nil {
		return nil, err
	}

	// Update the feeding schedule
	if err := schedule.ChangeSchedule(feedingTime, ConvertProtoFoodType(req.FoodType)); err != nil {
		return nil, err
	}

	// Save the feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return nil, err
	}

	return ConvertFeedingSchedule(schedule), nil
}

// CompleteFeedingSchedule implements the CompleteFeedingSchedule method of the FeedingScheduleService
func (s *FeedingScheduleServer) CompleteFeedingSchedule(ctx context.Context, req *zoo.CompleteFeedingScheduleRequest) (*zoo.Empty, error) {
	// Get the feeding schedule
	schedule, err := s.feedingScheduleRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Complete the feeding schedule
	if err := schedule.MarkCompleted(); err != nil {
		return nil, err
	}

	// Save the feeding schedule
	if err := s.feedingScheduleRepository.Save(schedule); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// GetDueFeedingSchedules implements the GetDueFeedingSchedules method of the FeedingScheduleService
func (s *FeedingScheduleServer) GetDueFeedingSchedules(ctx context.Context, req *zoo.Empty) (*zoo.GetFeedingSchedulesResponse, error) {
	// Get all feeding schedules
	schedules, err := s.feedingScheduleRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// Filter for due schedules
	var dueSchedules []*domain.FeedingSchedule
	for _, schedule := range schedules {
		if schedule.IsDue() && !schedule.Completed {
			dueSchedules = append(dueSchedules, schedule)
		}
	}

	// Convert to proto
	var protoSchedules []*zoo.FeedingSchedule
	for _, schedule := range dueSchedules {
		protoSchedules = append(protoSchedules, ConvertFeedingSchedule(schedule))
	}

	return &zoo.GetFeedingSchedulesResponse{
		FeedingSchedules: protoSchedules,
	}, nil
}

// GetFeedingSchedulesByAnimal implements the GetFeedingSchedulesByAnimal method of the FeedingScheduleService
func (s *FeedingScheduleServer) GetFeedingSchedulesByAnimal(ctx context.Context, req *zoo.GetFeedingSchedulesByAnimalRequest) (*zoo.GetFeedingSchedulesResponse, error) {
	// Check if the animal exists
	_, err := s.animalRepository.GetByID(req.AnimalId)
	if err != nil {
		return nil, err
	}

	// Get feeding schedules for the animal
	schedules, err := s.feedingScheduleRepository.GetByAnimalID(req.AnimalId)
	if err != nil {
		return nil, err
	}

	// Convert to proto
	var protoSchedules []*zoo.FeedingSchedule
	for _, schedule := range schedules {
		protoSchedules = append(protoSchedules, ConvertFeedingSchedule(schedule))
	}

	return &zoo.GetFeedingSchedulesResponse{
		FeedingSchedules: protoSchedules,
	}, nil
}
