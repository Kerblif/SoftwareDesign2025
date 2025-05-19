package servers

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"mini-dz-02/internal/proto/zoo"
	"mini-dz-02/pkg/zoo/domain"
)

// Convert domain.Gender to zoo.Gender
func ConvertGender(gender domain.Gender) zoo.Gender {
	switch gender {
	case domain.Male:
		return zoo.Gender_GENDER_MALE
	case domain.Female:
		return zoo.Gender_GENDER_FEMALE
	default:
		return zoo.Gender_GENDER_UNSPECIFIED
	}
}

// Convert domain.HealthStatus to zoo.HealthStatus
func ConvertHealthStatus(status domain.HealthStatus) zoo.HealthStatus {
	switch status {
	case domain.Healthy:
		return zoo.HealthStatus_HEALTH_STATUS_HEALTHY
	case domain.Sick:
		return zoo.HealthStatus_HEALTH_STATUS_SICK
	default:
		return zoo.HealthStatus_HEALTH_STATUS_UNSPECIFIED
	}
}

// Convert domain.EnclosureType to zoo.EnclosureType
func ConvertEnclosureType(enclosureType domain.EnclosureType) zoo.EnclosureType {
	switch enclosureType {
	case domain.Predator:
		return zoo.EnclosureType_ENCLOSURE_TYPE_PREDATOR
	case domain.Herbivore:
		return zoo.EnclosureType_ENCLOSURE_TYPE_HERBIVORE
	case domain.Aviary:
		return zoo.EnclosureType_ENCLOSURE_TYPE_AVIARY
	case domain.Aquarium:
		return zoo.EnclosureType_ENCLOSURE_TYPE_AQUARIUM
	case domain.Terrarium:
		return zoo.EnclosureType_ENCLOSURE_TYPE_TERRARIUM
	default:
		return zoo.EnclosureType_ENCLOSURE_TYPE_UNSPECIFIED
	}
}

// Convert domain.FoodType to zoo.FoodType
func ConvertFoodType(foodType domain.FoodType) zoo.FoodType {
	switch foodType {
	case domain.Meat:
		return zoo.FoodType_FOOD_TYPE_MEAT
	case domain.Vegetables:
		return zoo.FoodType_FOOD_TYPE_VEGETABLES
	case domain.Fruits:
		return zoo.FoodType_FOOD_TYPE_FRUITS
	case domain.Insects:
		return zoo.FoodType_FOOD_TYPE_INSECTS
	case domain.Seeds:
		return zoo.FoodType_FOOD_TYPE_SEEDS
	default:
		return zoo.FoodType_FOOD_TYPE_UNSPECIFIED
	}
}

// Convert zoo.Gender to domain.Gender
func ConvertProtoGender(gender zoo.Gender) domain.Gender {
	switch gender {
	case zoo.Gender_GENDER_MALE:
		return domain.Male
	case zoo.Gender_GENDER_FEMALE:
		return domain.Female
	default:
		return ""
	}
}

// Convert zoo.HealthStatus to domain.HealthStatus
func ConvertProtoHealthStatus(status zoo.HealthStatus) domain.HealthStatus {
	switch status {
	case zoo.HealthStatus_HEALTH_STATUS_HEALTHY:
		return domain.Healthy
	case zoo.HealthStatus_HEALTH_STATUS_SICK:
		return domain.Sick
	default:
		return ""
	}
}

// Convert zoo.EnclosureType to domain.EnclosureType
func ConvertProtoEnclosureType(enclosureType zoo.EnclosureType) domain.EnclosureType {
	switch enclosureType {
	case zoo.EnclosureType_ENCLOSURE_TYPE_PREDATOR:
		return domain.Predator
	case zoo.EnclosureType_ENCLOSURE_TYPE_HERBIVORE:
		return domain.Herbivore
	case zoo.EnclosureType_ENCLOSURE_TYPE_AVIARY:
		return domain.Aviary
	case zoo.EnclosureType_ENCLOSURE_TYPE_AQUARIUM:
		return domain.Aquarium
	case zoo.EnclosureType_ENCLOSURE_TYPE_TERRARIUM:
		return domain.Terrarium
	default:
		return ""
	}
}

// Convert zoo.FoodType to domain.FoodType
func ConvertProtoFoodType(foodType zoo.FoodType) domain.FoodType {
	switch foodType {
	case zoo.FoodType_FOOD_TYPE_MEAT:
		return domain.Meat
	case zoo.FoodType_FOOD_TYPE_VEGETABLES:
		return domain.Vegetables
	case zoo.FoodType_FOOD_TYPE_FRUITS:
		return domain.Fruits
	case zoo.FoodType_FOOD_TYPE_INSECTS:
		return domain.Insects
	case zoo.FoodType_FOOD_TYPE_SEEDS:
		return domain.Seeds
	default:
		return ""
	}
}

// Convert domain.Animal to zoo.Animal
func ConvertAnimal(animal *domain.Animal) *zoo.Animal {
	return &zoo.Animal{
		Id:           animal.ID,
		Species:      animal.Species.String(),
		Name:         animal.Name.String(),
		BirthDate:    timestamppb.New(animal.BirthDate.Time()),
		Gender:       ConvertGender(animal.Gender),
		FavoriteFood: ConvertFoodType(animal.FavoriteFood),
		HealthStatus: ConvertHealthStatus(animal.HealthStatus),
		EnclosureId:  animal.EnclosureID,
	}
}

// Convert domain.Enclosure to zoo.Enclosure
func ConvertEnclosure(enclosure *domain.Enclosure) *zoo.Enclosure {
	return &zoo.Enclosure{
		Id:               enclosure.ID,
		Type:             ConvertEnclosureType(enclosure.Type),
		Size:             int32(enclosure.Size.Value()),
		MaxCapacity:      int32(enclosure.MaxCapacity.Value()),
		CurrentAnimalIds: enclosure.CurrentAnimalIDs,
		IsClean:          enclosure.IsClean(),
	}
}

// Convert domain.FeedingSchedule to zoo.FeedingSchedule
func ConvertFeedingSchedule(schedule *domain.FeedingSchedule) *zoo.FeedingSchedule {
	return &zoo.FeedingSchedule{
		Id:          schedule.ID,
		AnimalId:    schedule.AnimalID,
		FeedingTime: timestamppb.New(schedule.FeedingTime.Time()),
		FoodType:    ConvertFoodType(schedule.FoodType),
		Completed:   schedule.Completed,
		IsDue:       schedule.IsDue(),
	}
}

// Convert map[string]int to map[string]int32
func ConvertMapStringInt(m map[string]int) map[string]int32 {
	result := make(map[string]int32)
	for k, v := range m {
		result[k] = int32(v)
	}
	return result
}

// Convert map[domain.Gender]int to map[string]int32
func ConvertMapGenderInt(m map[domain.Gender]int) map[string]int32 {
	result := make(map[string]int32)
	for k, v := range m {
		result[string(k)] = int32(v)
	}
	return result
}

// Convert map[domain.HealthStatus]int to map[string]int32
func ConvertMapHealthStatusInt(m map[domain.HealthStatus]int) map[string]int32 {
	result := make(map[string]int32)
	for k, v := range m {
		result[string(k)] = int32(v)
	}
	return result
}
