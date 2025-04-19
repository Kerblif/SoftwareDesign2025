package servers

import (
	"context"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/application"
)

// StatisticsServer implements the StatisticsService
type StatisticsServer struct {
	zoo.UnimplementedStatisticsServiceServer
	statisticsService *application.ZooStatisticsService
}

// NewStatisticsServer creates a new StatisticsServer
func NewStatisticsServer(
	statisticsService *application.ZooStatisticsService,
) *StatisticsServer {
	return &StatisticsServer{
		statisticsService: statisticsService,
	}
}

// GetZooStatistics implements the GetZooStatistics method of the StatisticsService
func (s *StatisticsServer) GetZooStatistics(ctx context.Context, req *zoo.Empty) (*zoo.ZooStatistics, error) {
	stats, err := s.statisticsService.GetZooStatistics()
	if err != nil {
		return nil, err
	}

	return &zoo.ZooStatistics{
		TotalAnimals:        int32(stats.TotalAnimals),
		HealthyAnimals:      int32(stats.HealthyAnimals),
		SickAnimals:         int32(stats.SickAnimals),
		TotalEnclosures:     int32(stats.TotalEnclosures),
		AvailableEnclosures: int32(stats.AvailableEnclosures),
		FullEnclosures:      int32(stats.FullEnclosures),
		EmptyEnclosures:     int32(stats.EmptyEnclosures),
		AnimalsBySpecies:    ConvertMapStringInt(stats.AnimalsBySpecies),
		AnimalsByGender:     ConvertMapGenderInt(stats.AnimalsByGender),
		AnimalsByEnclosure:  ConvertMapStringInt(stats.AnimalsByEnclosure),
	}, nil
}

// GetAnimalCountBySpecies implements the GetAnimalCountBySpecies method of the StatisticsService
func (s *StatisticsServer) GetAnimalCountBySpecies(ctx context.Context, req *zoo.Empty) (*zoo.AnimalCountBySpeciesResponse, error) {
	countBySpecies, err := s.statisticsService.GetAnimalCountBySpecies()
	if err != nil {
		return nil, err
	}

	return &zoo.AnimalCountBySpeciesResponse{
		CountBySpecies: ConvertMapStringInt(countBySpecies),
	}, nil
}

// GetEnclosureUtilization implements the GetEnclosureUtilization method of the StatisticsService
func (s *StatisticsServer) GetEnclosureUtilization(ctx context.Context, req *zoo.Empty) (*zoo.EnclosureUtilizationResponse, error) {
	utilization, err := s.statisticsService.GetEnclosureUtilization()
	if err != nil {
		return nil, err
	}

	return &zoo.EnclosureUtilizationResponse{
		Utilization: utilization,
	}, nil
}

// GetHealthStatusStatistics implements the GetHealthStatusStatistics method of the StatisticsService
func (s *StatisticsServer) GetHealthStatusStatistics(ctx context.Context, req *zoo.Empty) (*zoo.HealthStatusStatisticsResponse, error) {
	countByStatus, err := s.statisticsService.GetHealthStatusStatistics()
	if err != nil {
		return nil, err
	}

	return &zoo.HealthStatusStatisticsResponse{
		CountByStatus: ConvertMapHealthStatusInt(countByStatus),
	}, nil
}
