package servers

import (
	"context"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// EnclosureServer implements the EnclosureService
type EnclosureServer struct {
	zoo.UnimplementedEnclosureServiceServer
	enclosureManagementService *application.EnclosureManagementService
}

// NewEnclosureServer creates a new EnclosureServer
func NewEnclosureServer(
	enclosureManagementService *application.EnclosureManagementService,
) *EnclosureServer {
	return &EnclosureServer{
		enclosureManagementService: enclosureManagementService,
	}
}

// GetEnclosure implements the GetEnclosure method of the EnclosureService
func (s *EnclosureServer) GetEnclosure(ctx context.Context, req *zoo.GetEnclosureRequest) (*zoo.Enclosure, error) {
	enclosure, err := s.enclosureManagementService.GetEnclosureByID(req.Id)
	if err != nil {
		return nil, err
	}

	return ConvertEnclosure(enclosure), nil
}

// GetEnclosures implements the GetEnclosures method of the EnclosureService
func (s *EnclosureServer) GetEnclosures(ctx context.Context, req *zoo.Empty) (*zoo.GetEnclosuresResponse, error) {
	enclosures, err := s.enclosureManagementService.GetAllEnclosures()
	if err != nil {
		return nil, err
	}

	var protoEnclosures []*zoo.Enclosure
	for _, enclosure := range enclosures {
		protoEnclosures = append(protoEnclosures, ConvertEnclosure(enclosure))
	}

	return &zoo.GetEnclosuresResponse{
		Enclosures: protoEnclosures,
	}, nil
}

// CreateEnclosure implements the CreateEnclosure method of the EnclosureService
func (s *EnclosureServer) CreateEnclosure(ctx context.Context, req *zoo.CreateEnclosureRequest) (*zoo.Enclosure, error) {
	// Create value objects
	size, err := domain.NewEnclosureSize(int(req.Size))
	if err != nil {
		return nil, err
	}

	capacity, err := domain.NewCapacity(int(req.MaxCapacity))
	if err != nil {
		return nil, err
	}

	// Create the enclosure using the service
	enclosure, err := s.enclosureManagementService.CreateEnclosure(
		req.Id,
		ConvertProtoEnclosureType(req.Type),
		size,
		capacity,
	)
	if err != nil {
		return nil, err
	}

	return ConvertEnclosure(enclosure), nil
}

// DeleteEnclosure implements the DeleteEnclosure method of the EnclosureService
func (s *EnclosureServer) DeleteEnclosure(ctx context.Context, req *zoo.DeleteEnclosureRequest) (*zoo.Empty, error) {
	// Delete the enclosure using the service
	if err := s.enclosureManagementService.DeleteEnclosure(req.Id); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// CleanEnclosure implements the CleanEnclosure method of the EnclosureService
func (s *EnclosureServer) CleanEnclosure(ctx context.Context, req *zoo.CleanEnclosureRequest) (*zoo.Empty, error) {
	// Clean the enclosure using the service
	if err := s.enclosureManagementService.CleanEnclosure(req.Id); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// GetAnimalsInEnclosure implements the GetAnimalsInEnclosure method of the EnclosureService
func (s *EnclosureServer) GetAnimalsInEnclosure(ctx context.Context, req *zoo.GetAnimalsInEnclosureRequest) (*zoo.GetAnimalsInEnclosureResponse, error) {
	// Get the animals in the enclosure using the service
	animals, err := s.enclosureManagementService.GetAnimalsInEnclosure(req.Id)
	if err != nil {
		return nil, err
	}

	// Convert to proto
	var protoAnimals []*zoo.Animal
	for _, animal := range animals {
		protoAnimals = append(protoAnimals, ConvertAnimal(animal))
	}

	return &zoo.GetAnimalsInEnclosureResponse{
		Animals: protoAnimals,
	}, nil
}
