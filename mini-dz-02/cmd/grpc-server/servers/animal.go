package servers

import (
	"context"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/application"
	"mini-dz-02/pkg/zoo/domain"
)

// AnimalServer implements the AnimalService
type AnimalServer struct {
	zoo.UnimplementedAnimalServiceServer
	animalRepository    domain.AnimalRepository
	enclosureRepository domain.EnclosureRepository
	transferService     *application.AnimalTransferService
}

// NewAnimalServer creates a new AnimalServer
func NewAnimalServer(
	animalRepository domain.AnimalRepository,
	enclosureRepository domain.EnclosureRepository,
	transferService *application.AnimalTransferService,
) *AnimalServer {
	return &AnimalServer{
		animalRepository:    animalRepository,
		enclosureRepository: enclosureRepository,
		transferService:     transferService,
	}
}

// GetAnimal implements the GetAnimal method of the AnimalService
func (s *AnimalServer) GetAnimal(ctx context.Context, req *zoo.GetAnimalRequest) (*zoo.Animal, error) {
	animal, err := s.animalRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return ConvertAnimal(animal), nil
}

// GetAnimals implements the GetAnimals method of the AnimalService
func (s *AnimalServer) GetAnimals(ctx context.Context, req *zoo.Empty) (*zoo.GetAnimalsResponse, error) {
	animals, err := s.animalRepository.GetAll()
	if err != nil {
		return nil, err
	}

	var protoAnimals []*zoo.Animal
	for _, animal := range animals {
		protoAnimals = append(protoAnimals, ConvertAnimal(animal))
	}

	return &zoo.GetAnimalsResponse{
		Animals: protoAnimals,
	}, nil
}

// CreateAnimal implements the CreateAnimal method of the AnimalService
func (s *AnimalServer) CreateAnimal(ctx context.Context, req *zoo.CreateAnimalRequest) (*zoo.Animal, error) {
	// Create value objects
	species, err := domain.NewSpecies(req.Species)
	if err != nil {
		return nil, err
	}

	name, err := domain.NewAnimalName(req.Name)
	if err != nil {
		return nil, err
	}

	birthDate, err := domain.NewBirthDate(req.BirthDate.AsTime())
	if err != nil {
		return nil, err
	}

	// Create the animal
	animal, err := domain.NewAnimal(
		req.Id,
		species,
		name,
		birthDate,
		ConvertProtoGender(req.Gender),
		ConvertProtoFoodType(req.FavoriteFood),
		ConvertProtoHealthStatus(req.HealthStatus),
	)
	if err != nil {
		return nil, err
	}

	// Save the animal
	if err := s.animalRepository.Save(animal); err != nil {
		return nil, err
	}

	return ConvertAnimal(animal), nil
}

// DeleteAnimal implements the DeleteAnimal method of the AnimalService
func (s *AnimalServer) DeleteAnimal(ctx context.Context, req *zoo.DeleteAnimalRequest) (*zoo.Empty, error) {
	// Check if the animal exists
	animal, err := s.animalRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// If the animal is in an enclosure, remove it first
	if animal.IsInEnclosure() {
		enclosure, err := s.enclosureRepository.GetByID(animal.EnclosureID)
		if err != nil {
			return nil, err
		}

		if err := enclosure.RemoveAnimal(animal.ID); err != nil {
			return nil, err
		}

		if err := s.enclosureRepository.Save(enclosure); err != nil {
			return nil, err
		}
	}

	// Delete the animal
	if err := s.animalRepository.Delete(req.Id); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// TransferAnimal implements the TransferAnimal method of the AnimalService
func (s *AnimalServer) TransferAnimal(ctx context.Context, req *zoo.TransferAnimalRequest) (*zoo.Empty, error) {
	// Transfer the animal
	if err := s.transferService.TransferAnimal(req.Id, req.EnclosureId); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}

// TreatAnimal implements the TreatAnimal method of the AnimalService
func (s *AnimalServer) TreatAnimal(ctx context.Context, req *zoo.TreatAnimalRequest) (*zoo.Empty, error) {
	// Get the animal
	animal, err := s.animalRepository.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// Treat the animal
	if err := animal.Treat(); err != nil {
		return nil, err
	}

	// Save the animal
	if err := s.animalRepository.Save(animal); err != nil {
		return nil, err
	}

	return &zoo.Empty{}, nil
}
