package servers

import (
	"context"
	"fmt"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/domain"
)

// EnclosureServer implements the EnclosureService
type EnclosureServer struct {
	zoo.UnimplementedEnclosureServiceServer
	enclosureRepository domain.EnclosureRepository
	animalRepository    domain.AnimalRepository
}

// NewEnclosureServer creates a new EnclosureServer
func NewEnclosureServer(
	enclosureRepository domain.EnclosureRepository,
	animalRepository domain.AnimalRepository,
) *EnclosureServer {
	return &EnclosureServer{
		enclosureRepository: enclosureRepository,
		animalRepository:    animalRepository,
	}
}

// GetEnclosure implements the GetEnclosure method of the EnclosureService
func (s *EnclosureServer) GetEnclosure(ctx context.Context, req *zoo.GetEnclosureRequest) (*zoo.Enclosure, error) {
	enclosure, err := s.enclosureRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return ConvertEnclosure(enclosure), nil
}

// GetEnclosures implements the GetEnclosures method of the EnclosureService
func (s *EnclosureServer) GetEnclosures(ctx context.Context, req *zoo.Empty) (*zoo.GetEnclosuresResponse, error) {
	enclosures, err := s.enclosureRepository.GetAll()
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

	// Create the enclosure
	enclosure, err := domain.NewEnclosure(
		req.Id,
		ConvertProtoEnclosureType(req.Type),
		size,
		capacity,
	)
	if err != nil {
		return nil, err
	}

	// Save the enclosure
	if err := s.enclosureRepository.Save(enclosure); err != nil {
		return nil, err
	}

	return ConvertEnclosure(enclosure), nil
}

// DeleteEnclosure implements the DeleteEnclosure method of the EnclosureService
func (s *EnclosureServer) DeleteEnclosure(ctx context.Context, req *zoo.DeleteEnclosureRequest) (*zoo.Empty, error) {
	// Check if the enclosure exists
	enclosure, err := s.enclosureRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Check if the enclosure is empty
	if len(enclosure.CurrentAnimalIDs) > 0 {
		return nil, fmt.Errorf("cannot delete enclosure with animals")
	}

	// Delete the enclosure
	if err := s.enclosureRepository.Delete(req.Id); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// CleanEnclosure implements the CleanEnclosure method of the EnclosureService
func (s *EnclosureServer) CleanEnclosure(ctx context.Context, req *zoo.CleanEnclosureRequest) (*zoo.Empty, error) {
	// Get the enclosure
	enclosure, err := s.enclosureRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Clean the enclosure
	if err := enclosure.Clean(); err != nil {
		return nil, err
	}

	// Save the enclosure
	if err := s.enclosureRepository.Save(enclosure); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// GetAnimalsInEnclosure implements the GetAnimalsInEnclosure method of the EnclosureService
func (s *EnclosureServer) GetAnimalsInEnclosure(ctx context.Context, req *zoo.GetAnimalsInEnclosureRequest) (*zoo.GetAnimalsInEnclosureResponse, error) {
	// Get the enclosure
	enclosure, err := s.enclosureRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Get the animals in the enclosure
	var protoAnimals []*zoo.Animal
	for _, animalID := range enclosure.CurrentAnimalIDs {
		animal, err := s.animalRepository.GetByID(animalID)
		if err != nil {
			return nil, err
		}

		protoAnimals = append(protoAnimals, ConvertAnimal(animal))
	}

	return &zoo.GetAnimalsInEnclosureResponse{
		Animals: protoAnimals,
	}, nil
}
